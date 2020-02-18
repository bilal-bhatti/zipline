package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

type Zipline struct {
	packets   []*packet
	templates map[string]*template
	renderer  *renderer
	provider  *provider
	fset      *token.FileSet
}

func NewZipline() *Zipline {
	return &Zipline{
		packets:   []*packet{},
		templates: make(map[string]*template),
		fset:      token.NewFileSet(),
	}
}

func (z *Zipline) Start() {
	pkgs := load()
	scanner := scanner{pkgs: pkgs}
	z.templates, z.packets = scanner.scan()

	z.provider = newProvider(pkgs)
	z.renderer = newRenderer(z.templates, z.provider)

	for _, packet := range z.packets {
		z.prepare(packet)

		z.renderer.complete(packet)

		root, err := goSrcRoot()
		if err != nil {
			panic(err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		od := strings.TrimPrefix(packet.pkg.PkgPath, root)
		// log.Println("cwd", cwd)
		// log.Println("pkg path", packet.pkg.PkgPath)
		// log.Println("source root", root)
		// log.Println("output dir", od)

		out := path.Join(cwd, od, "bindings_gen.go")

		f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			panic(err)
		}
		// z.renderer.print(os.Stdout, false) // TODO: in case of error dump out for debug
		z.renderer.print(f, true)
		f.Close()

		bs, err := ioutil.ReadFile(out)
		if err != nil {
			panic(err)
		}
		opt := imports.Options{
			Comments:   true,
			FormatOnly: false,
		}
		bs, err = imports.Process(out, bs, &opt)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(out, bs, os.ModePerm)
		if err != nil {
			panic(err)
		}

		log.Println("Wrote bindings to", out)
	}

	swagger := newSwagger()
	swagger.generate(z.packets)
}

func (z *Zipline) prepare(packet *packet) {
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
					log.Println("Somethinger other than call to zipline found")
					continue
				}

				id, ok := sel.X.(*ast.Ident)
				if !ok {
					continue
				}

				ido := z.provider.qualifiedIdentObject(id)

				// ensure var type is ZiplineTemplate
				if ido != nil && strings.HasSuffix(ido.Type().String(), ZiplineTemplate) {
					// generate function body first
					binding := z.processStatement(packet.pkg, expStmt)

					// rewrite ast to replace zipline spec
					expType.Args[i] = newCallExpression(binding, expType.Args[i])
					packet.bindings = append(packet.bindings, binding)
				}
			}
		default:
			log.Println("Unhandled expression type", reflect.TypeOf(expType))
		}
	}
}

func (z *Zipline) processStatement(pkg *packages.Package, stmt *ast.ExprStmt) *binding {
	binding := parseSpec(pkg, stmt)
	z.renderer.render(pkg, binding)
	return binding
}

func newCallExpression(binding *binding, arg ast.Expr) *ast.CallExpr {
	ce := &ast.CallExpr{
		Fun: &ast.Ident{
			Name: binding.id() + "HandlerFunc", // TODO: use the template return type
		},
	}
	return ce
}

func parseSpec(pkg *packages.Package, spec *ast.ExprStmt) *binding {
	call, ok := spec.X.(*ast.CallExpr)
	if !ok {
		panic("spec invalid")
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		panic("spec invalid")
	}

	path := strings.Trim(call.Args[0].(*ast.BasicLit).Value, "\"")

	binding := &binding{
		template:       sel.Sel.Name,
		path:           path,
		paramTemplates: []string{},
	}

	zipline, ok := call.Args[1].(*ast.CallExpr)
	if !ok {
		panic("invalid expression")
	}

	switch handler := zipline.Args[0].(type) {
	case *ast.SelectorExpr:
		binding.handler = newHandlerInfoFromSelectorExpr(pkg, handler)
	case *ast.Ident:
		binding.handler = newHandlerInfoFromIdent(pkg, handler)
	default:
		panic("unsupported expression")
	}

	// parse additional parameters, if any
	for i := 1; i < len(zipline.Args); i++ {
		arg := zipline.Args[i]

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

	return binding
}

func newHandlerInfoFromSelectorExpr(pkg *packages.Package, handler *ast.SelectorExpr) *handlerInfo {
	hi := newHandlerInfoFromIdent(pkg, handler.Sel)

	switch xt := handler.X.(type) {
	case *ast.Ident:
		// same package
		obj := qualifiedIdentObject(pkg.TypesInfo, xt)
		if obj != nil {
			if _, ok := obj.Type().(*types.Basic); !ok {
				hi.x = newTypeToken("", obj.Type().String(), "")
			}
		}
	case *ast.SelectorExpr:
		// different package
		// xt.X = package
		// xt.Sel = type
		obj := qualifiedIdentObject(pkg.TypesInfo, xt.Sel)
		if obj != nil {
			hi.x = newTypeToken("", obj.Type().String(), "")
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
						hi.x = newTypeToken("", obj.Type().String(), "")
					}
				case *ast.Ident:
					// same package
					obj := qualifiedIdentObject(pkg.TypesInfo, newExpType)
					if obj != nil {
						if _, ok := obj.Type().(*types.Basic); !ok {
							hi.x = newTypeToken("", obj.Type().String(), "")
						}
					}
				}
			}
		}
	default:
		msg := &bytes.Buffer{}
		printNode(msg, handler)
		panic(fmt.Sprintf("invalid zipline template parameter %s", msg.String()))
	}

	return hi
}

func newHandlerInfoFromIdent(pkg *packages.Package, handler *ast.Ident) *handlerInfo {
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

	hi := &handlerInfo{
		id:  string(id.Bytes()),
		sel: handler.String(),
		pkg: obj.Pkg().Path(),
	}

	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		tt := newTypeToken(pkg.Name, p.Type().String(), p.Name())
		tt.varType = p
		hi.params = append(hi.params, tt)
	}

	for i := 0; i < sig.Results().Len(); i++ {
		r := sig.Results().At(i)
		tt := newTypeToken(pkg.Name, r.Type().String(), r.Name())
		tt.varType = r
		hi.returns = append(hi.returns, tt)
	}

	return hi
}
