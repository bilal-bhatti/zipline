package internal

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"strings"

	"github.com/go-toolsmith/astcopy"
	"golang.org/x/tools/go/packages"
)

type packet struct {
	pkg      *packages.Package
	bindings *ast.FuncDecl
}

type Zipline struct {
	packets  []*packet
	renderer *renderer
}

func NewZipline() *Zipline {
	return &Zipline{
		packets:  []*packet{},
		renderer: newRenderer(),
	}
}

func (z *Zipline) Start() {
	pkgs := load()

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, d := range file.Decls {
				switch t := d.(type) {
				case *ast.FuncDecl:
					if isBindingSpecNode(pkg.TypesInfo, t) {
						z.packets = append(z.packets, &packet{
							pkg:      pkg,
							bindings: t,
						})
					}
				}
			}
		}

		for k, v := range pkg.TypesInfo.Defs {
			if strings.Contains(k.Name, "Provide") {
				log.Println("k,v", k, v)
			}
			if strings.Contains(k.Name, "Response") {
				log.Println("k,v", k, v, reflect.TypeOf(v.Type()))

			}
		}
	}

	for _, packet := range z.packets {
		// 	wr := astutil.Apply(packet.Bindings, func(c *astutil.Cursor) bool {
		// 		log.Println("pre", c.Node())

		// 		return true
		// 	}, func(c *astutil.Cursor) bool {
		// 		return true
		// 	})

		// 	fset := token.NewFileSet()
		// 	var buf bytes.Buffer
		// 	printer.Fprint(&buf, fset, wr)
		// 	log.Println("body", string(buf.Bytes()))
		// rewrite(packet.Pkg.TypesInfo, packet.Bindings)
		z.prepare(packet)
		// parseSpecExpression(packet)
	}
}

func (z *Zipline) prepare(packet *packet) {
	dfunc := packet.bindings
	// body := bytes.Buffer{}

	for _, stmt := range dfunc.Body.List {
		switch stmtType := stmt.(type) {
		case *ast.ExprStmt:
			if isBindingSpecNode(packet.pkg.TypesInfo, stmtType) {
				print("binding spec", stmtType)
				switch expType := stmtType.X.(type) {
				case *ast.CallExpr:

					for i := 0; i < len(expType.Args); i++ {
						arg := expType.Args[i]
						// log.Println("arg", arg, reflect.TypeOf(arg))

						if call, ok := arg.(*ast.CallExpr); ok {
							if isBindingSpecNode(packet.pkg.TypesInfo, call) {
								// log.Println("call fun", reflect.TypeOf(call.Fun))
								if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
									if id, ok := sel.X.(*ast.Ident); ok {
										if id.Name == "zipline" {

											// generate(packet.pkg.TypesInfo, sel)
											z.processStatement(packet.pkg.TypesInfo, stmtType)
											// expType.Args[i] = ast.NewIdent(generate(info, sel))

											// ContactsServiceCreateHandlerFunc
											handler := astcopy.Expr(expType.Args[i])

											expType.Args[i] = newCallExpression(handler)
										}
									}
								} else {
									log.Println("Somethinger other than call to zipling found")
								}
							}
						}
					}
				default:
					log.Println("Unhandled exp type", reflect.TypeOf(expType))
				}
			}
		}
	}

	fset := token.NewFileSet()
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, packet.bindings)
	log.Println("func body \n", string(buf.Bytes()))
	z.renderer.complete()
	z.renderer.print(true)
}

func (z *Zipline) processStatement(info *types.Info, stmt *ast.ExprStmt) {
	binding := parseSpec(info, stmt)

	z.renderer.render(binding)
	printBinding(binding)
}

func newCallExpression(arg ast.Expr) *ast.CallExpr {
	ce := &ast.CallExpr{
		Fun: &ast.Ident{
			Name: "ContactsServiceCreateHandlerFunc",
		},
		Args: arg.(*ast.CallExpr).Args,
	}
	return ce
}

func printBinding(b *binding) {
	log.Println(" *** ")
	log.Println("path", b.path)
	log.Println("method", b.method)
	for _, p := range b.handler.params {
		log.Println("p", p.Name, p.Type)
	}
	for _, r := range b.handler.returns {
		log.Println("r", r.Name, r.Type)
	}
}

func parseSpec(info *types.Info, spec *ast.ExprStmt) *binding {
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

	// ast.Print(token.NewFileSet(), call.Args[1])

	zipline, ok := call.Args[1].(*ast.CallExpr)
	if !ok {
		panic("invalid expression")
	}

	switch handler := zipline.Args[0].(type) {
	case *ast.SelectorExpr:
		binding.handler = newFromSelectorExpr(info, handler)
	case *ast.Ident:
		binding.handler = newFromIdent(info, handler)
	}

	return binding
}

func newFromSelectorExpr(info *types.Info, handler *ast.SelectorExpr) *handlerInfo {
	hi := &handlerInfo{}
	// switch xt := handler.X.(type) {
	// case *ast.CallExpr:
	// 	xo := qualifiedIdentObject(info, xt.Fun)

	// 	if xv, ok := xo.(*types.Func); ok {
	// 		log.Println("xv", xv.Type())
	// 	}
	// case *ast.Ident:
	// 	xo := qualifiedIdentObject(info, handler.X)

	// 	if xv, ok := xo.(*types.Var); ok {
	// 		log.Println("xv", xv.Type())
	// 	}
	// default:
	// 	panic(fmt.Sprintf("Unhandled zipline spec expresion of type %v", reflect.TypeOf(xt)))
	// }

	obj := qualifiedIdentObject(info, handler.Sel)

	sig := obj.Type().(*types.Signature)

	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		hi.params = append(hi.params, &BindingVar{
			Name: p.Name(),
			Type: p.Type().String(),
		})
	}

	for i := 0; i < sig.Results().Len(); i++ {
		r := sig.Results().At(i)
		hi.returns = append(hi.returns, &BindingVar{
			Name: r.Name(),
			Type: r.Type().String(),
		})
	}

	return hi
}

func newFromIdent(info *types.Info, handler *ast.Ident) *handlerInfo {
	// log.Println(handler)

	hi := &handlerInfo{}

	obj := qualifiedIdentObject(info, handler)

	// if xv, ok := obj.(*types.Func); ok {
	// 	log.Println("xv", xv.Type())
	// }

	sig := obj.Type().(*types.Signature)

	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		hi.params = append(hi.params, &BindingVar{
			Name: p.Name(),
			Type: p.Type().String(),
		})
	}

	for i := 0; i < sig.Results().Len(); i++ {
		r := sig.Results().At(i)
		hi.returns = append(hi.returns, &BindingVar{
			Name: r.Name(),
			Type: r.Type().String(),
		})
	}

	return hi
}
