package internal

import (
	"bytes"
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

type packet struct {
	pkg      *packages.Package
	bindings *ast.FuncDecl
}

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
}

func (z *Zipline) prepare(packet *packet) {
	dfunc := packet.bindings

	// TODO: clean this up. use ast.Inspect to simplify this mess
	for _, stmt := range dfunc.Body.List {
		switch stmtType := stmt.(type) {
		case *ast.ExprStmt:
			// does the statement contain zipline binding
			if isBindingSpecNode(packet.pkg.TypesInfo, stmtType) {
				switch expType := stmtType.X.(type) {
				case *ast.CallExpr:
					for i := 0; i < len(expType.Args); i++ {
						arg := expType.Args[i]
						if call, ok := arg.(*ast.CallExpr); ok {
							if isBindingSpecNode(packet.pkg.TypesInfo, call) {
								// actual call to ZiplineTemplate
								if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
									if id, ok := sel.X.(*ast.Ident); ok {
										ido := z.provider.qualifiedIdentObject(id)

										// ensure var type is ZiplineTemplate
										if ido != nil && strings.HasSuffix(ido.Type().String(), ZiplineTemplate) {
											// generate function body first
											binding := z.processStatement(packet.pkg, stmtType)

											// rewrite ast to replace zipline spec
											expType.Args[i] = newCallExpression(binding, expType.Args[i])
										}
									}
								} else {
									log.Println("Somethinger other than call to zipline found")
								}
							}
						}
					}
				default:
					log.Println("Unhandled expression type", reflect.TypeOf(expType))
				}
			}
		}
	}
}

func (z *Zipline) processStatement(pkg *packages.Package, stmt *ast.ExprStmt) *binding {
	binding := parseSpec(pkg, stmt)
	z.renderer.render(binding)
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
		method: sel.Sel.Name,
		path:   path,
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
	}

	return binding
}

func newHandlerInfoFromSelectorExpr(pkg *packages.Package, handler *ast.SelectorExpr) *handlerInfo {
	hi := newHandlerInfoFromIdent(pkg, handler.Sel)

	xid, ok := handler.X.(*ast.Ident)
	if !ok {
		panic("Zipline must be a function selector expression (X.Y) where X is a service struct and Y is a method")
	}

	obj := qualifiedIdentObject(pkg.TypesInfo, xid)

	hi.x = newTypeToken("", obj.Type().String(), "")

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
	}

	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		hi.params = append(hi.params, newTypeToken(pkg.Name, p.Type().String(), p.Name()))
	}

	for i := 0; i < sig.Results().Len(); i++ {
		r := sig.Results().At(i)
		hi.returns = append(hi.returns, newTypeToken(pkg.Name, r.Type().String(), r.Name()))
	}

	return hi
}
