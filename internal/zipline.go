package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

type Zipline struct {
	packets   []*packet
	templates map[string]*template
	typeSpecs map[string]*typeSpecWithPkg
	renderer  *renderer
	provider  *provider
}

func NewZipline() *Zipline {
	return &Zipline{
		packets:   []*packet{},
		templates: make(map[string]*template),
		typeSpecs: make(map[string]*typeSpecWithPkg),
	}
}

func (z *Zipline) Start(pkgPaths []string) error {
	cw, err := os.Getwd()

	// log current directory
	log.Println(cw)

	pkgs, err := load(pkgPaths)
	if err != nil {
		return err
	}
	scanner := scanner{pkgs: pkgs}
	z.typeSpecs, z.templates, z.packets = scanner.scan()

	z.provider = newProvider(pkgs)
	z.renderer = newRenderer(z.templates, z.provider)

	for _, packet := range z.packets {
		if err := z.prepare(packet); err != nil {
			return err
		}

		z.renderer.complete(packet)

		output := bytes.NewBuffer(make([]byte, 0))
		err := z.renderer.print(output, true)
		if err != nil {
			return err
		}

		root, err := goSrcRoot()
		if err != nil {
			return err
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		od := strings.TrimPrefix(packet.pkg.PkgPath, root)
		trace("- calculating output package location -")
		trace("cwd: %s", cwd)
		trace("package path: %s", packet.pkg.PkgPath)
		trace("source root: %s", root)
		trace("output dir: %s", od)

		out := path.Join(cwd, od, "bindings_gen.go")

		f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		// z.renderer.print(os.Stdout, false) // TODO: in case of error dump out for debug

		if err != nil {

		}
		_, err = f.Write(output.Bytes())
		if err != nil {
			f.Close()
			return err
		}

		f.Close()

		bs, err := ioutil.ReadFile(out)
		if err != nil {
			return err
		}
		opt := imports.Options{
			Comments:   true,
			FormatOnly: false,
		}
		bs, err = imports.Process(out, bs, &opt)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(out, bs, os.ModePerm)
		if err != nil {
			return err
		}

		log.Printf("wrote bindings to .%s\n", strings.TrimPrefix(out, cwd))

		// reset code buffers
		z.renderer.body.buf.Reset()
		z.renderer.preamble.buf.Reset()
	}

	swagger, err := newSwagger(z.typeSpecs)
	if err != nil {
		return err
	}

	err = swagger.generate(z.packets)
	if err != nil {
		return err
	}

	return nil
}

func (z *Zipline) prepare(packet *packet) error {
	dfunc := packet.funcDecl

	for _, stmt := range dfunc.Body.List {
		// looking for expression statements
		expStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			// keep going
			continue
		}

		// does the statement contain zipline binding
		if !isZiplineNode(packet.pkg.TypesInfo, expStmt) {
			// if not
			continue
		}

		switch expType := expStmt.X.(type) {
		case *ast.CallExpr:
			// call expression, let's examine the arguments
			for i := 0; i < len(expType.Args); i++ {
				arg := expType.Args[i]

				// is argument a call expression
				call, ok := arg.(*ast.CallExpr)
				if !ok {
					// if not
					continue
				}

				// is call expression a zipline call expression
				if !isZiplineNode(packet.pkg.TypesInfo, call) {
					continue
				}

				// actual call to ZiplineTemplate, i.e. ZiplineTemplate.TemplateFunc
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					// TODO: fix this
					// handler/zipline call is being wrapped
					// should handle this properly

					return newErrorForStmt("failed to parse zipline expression, call likely wrapped", stmt)
				}

				id, ok := sel.X.(*ast.Ident)
				if !ok {
					continue
				}

				ido := z.provider.qualifiedIdentObject(id)

				// ensure var type is ZiplineTemplate
				if ido != nil && strings.HasSuffix(ido.Type().String(), ZiplineTemplate) {
					// generate function body first
					binding, err := z.processStatement(packet.pkg, expStmt)
					if err != nil {
						return err
					}

					// rewrite ast to replace zipline spec
					expType.Args[i] = newCallExpression(binding, expType.Args[i])
					packet.bindings = append(packet.bindings, binding)
				}
			}
		default:
			return newErrorForStmt(fmt.Sprintf("unhandled expression type %v", reflect.TypeOf(expType)), stmt)
		}
	}

	return nil
}

func (z *Zipline) processStatement(pkg *packages.Package, stmt *ast.ExprStmt) (*binding, error) {
	binding, err := parseSpec(pkg, stmt)
	if err != nil {
		return nil, err
	}

	if err := z.renderer.render(pkg, binding); err != nil {
		return nil, err
	}
	return binding, nil
}

func newCallExpression(binding *binding, arg ast.Expr) *ast.CallExpr {
	ce := &ast.CallExpr{
		Fun: &ast.Ident{
			Name: binding.id() + "HandlerFunc", // TODO: use the template return type
		},
	}

	if len(binding.boundParams) > 0 {
		args := make([]ast.Expr, len(binding.boundParams))
		for i, t := range binding.boundParams {
			args[i] = ast.NewIdent(t.varName())
		}
		ce.Args = args
	}

	return ce
}

func parseSpec(pkg *packages.Package, spec *ast.ExprStmt) (*binding, error) {
	call, ok := spec.X.(*ast.CallExpr)
	if !ok {
		return nil, errors.New("spec invalid")
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, errors.New("spec invalid")
	}

	path := strings.Trim(call.Args[0].(*ast.BasicLit).Value, "\"")

	zipline, ok := call.Args[1].(*ast.CallExpr)
	if !ok {
		return nil, errors.New("invalid expression")
	}

	binding := &binding{
		template:       sel.Sel.Name,
		path:           path,
		paramTemplates: []string{},
	}

	// capture zipline spec for error reporting
	zsb := newBuffer()
	printer.Fprint(zsb.buf, pkg.Fset, zipline)
	binding.spec = string(zsb.buf.Bytes())

	switch handler := zipline.Args[0].(type) {
	case *ast.SelectorExpr:
		handle, err := newHandlerInfoFromSelectorExpr(pkg, handler)
		if err != nil {
			return nil, err
		}

		binding.handler = handle
	case *ast.Ident:
		handlerInfo, err := newHandlerInfoFromIdent(pkg, handler)
		if err != nil {
			return nil, err
		}
		binding.handler = handlerInfo
	default:
		return nil, errors.New("unsupported expression")
	}

	// parse additional parameters, if any
	for i := 1; i < len(zipline.Args); i++ {
		arg := zipline.Args[i]

		if ident, ok := arg.(*ast.Ident); ok && ident.Obj.Kind == ast.Var {
			// this is a var being passed down through the template
			xo := qualifiedIdentObject(pkg.TypesInfo, ident)
			binding.boundParams = append(binding.boundParams, newTypeToken("", xo.Type().String(), ident.Name))
			continue
		}

		expr, ok := arg.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		if x, ok := expr.X.(*ast.Ident); ok {
			xo := qualifiedIdentObject(pkg.TypesInfo, x)

			// if receiver is ZiplineTemplate
			if strings.HasSuffix(xo.Type().String(), ZiplineTemplate) {
				binding.paramTemplates = append(binding.paramTemplates, expr.Sel.Name)
			}
		}
	}

	return binding, nil
}

func newHandlerInfoFromSelectorExpr(pkg *packages.Package, handler *ast.SelectorExpr) (*handlerInfo, error) {
	handlerInfo, err := newHandlerInfoFromIdent(pkg, handler.Sel)
	if err != nil {
		return nil, err
	}

	switch xt := handler.X.(type) {
	case *ast.Ident:
		// same package
		obj := qualifiedIdentObject(pkg.TypesInfo, xt)
		if obj != nil {
			if _, ok := obj.Type().(*types.Basic); !ok {
				handlerInfo.x = newTypeToken("", obj.Type().String(), "")
			}
		}
	case *ast.SelectorExpr:
		// different package
		// xt.X = package
		// xt.Sel = type
		obj := qualifiedIdentObject(pkg.TypesInfo, xt.Sel)
		if obj != nil {
			handlerInfo.x = newTypeToken("", obj.Type().String(), "")
		}
	case *ast.CallExpr:
		// if it's a a new call
		if id, ok := xt.Fun.(*ast.Ident); ok {
			if id.String() == "new" {
				switch newExpType := xt.Args[0].(type) {
				case *ast.SelectorExpr:
					// different package
					obj := qualifiedIdentObject(pkg.TypesInfo, newExpType.Sel)
					if obj != nil {
						handlerInfo.x = newTypeToken("", obj.Type().String(), "")
					}
				case *ast.Ident:
					// same package
					obj := qualifiedIdentObject(pkg.TypesInfo, newExpType)
					if obj != nil {
						if _, ok := obj.Type().(*types.Basic); !ok {
							handlerInfo.x = newTypeToken("", obj.Type().String(), "")
						}
					}
				}
			}
		}
	default:
		msg := newBuffer()
		printNode(msg, handler)
		return nil, errors.New(fmt.Sprintf("invalid zipline template parameter %s", msg.buf.String()))
	}

	return handlerInfo, nil
}

func newHandlerInfoFromIdent(pkg *packages.Package, handler *ast.Ident) (*handlerInfo, error) {
	obj := qualifiedIdentObject(pkg.TypesInfo, handler)

	sig := obj.Type().(*types.Signature)

	var id bytes.Buffer

	// id.WriteString(strings.Title(obj.Pkg().Name()))
	if sig.Recv() != nil {
		rcvr := sig.Recv().Type().String()
		idx := strings.LastIndex(rcvr, ".")

		if idx > 0 {
			rcvr = strings.Trim(rcvr[idx:len(rcvr)], ".")
		}
		id.WriteString(rcvr)
	}
	id.WriteString(obj.Name())

	pos := pkg.Fset.PositionFor(obj.Pos(), true)
	comments, err := getComments(pos)
	if err != nil {
		// let's not fail on comments but log the error
		log.Println("failed to extract comments", err.Error())
	}

	hi := &handlerInfo{
		comments:  comments,
		id:        string(id.Bytes()),
		sel:       handler.String(),
		pkg:       obj.Pkg().Path(),
		signature: sig,
	}

	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		tt := newTypeToken(pkg.Name, p.Type().String(), p.Name())
		tt.varType = p
		hi.params = append(hi.params, tt)
	}

	for i := 0; i < sig.Results().Len(); i++ {
		r := sig.Results().At(i)
		if _, ok := r.Type().(*types.Slice); ok {
			return nil, newErrorForSliceVar("return type should not be a slice", obj)
		}
		tt := newTypeToken(pkg.Name, r.Type().String(), r.Name())
		tt.varType = r
		hi.returns = append(hi.returns, tt)
	}

	return hi, nil
}
