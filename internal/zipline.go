package internal

import (
	"bytes"
	"go/ast"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/go-toolsmith/astcopy"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

type packet struct {
	pkg      *packages.Package
	bindings *ast.FuncDecl
}

type Zipline struct {
	packets  []*packet
	renderer *renderer
	provider *provider
}

func NewZipline() *Zipline {
	return &Zipline{
		packets: []*packet{},
	}
}

func (z *Zipline) Start() {
	pkgs := load()
	z.provider = newProvider(pkgs)
	z.renderer = newRenderer(z.provider)

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
	}

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

	for _, stmt := range dfunc.Body.List {
		switch stmtType := stmt.(type) {
		case *ast.ExprStmt:
			if isBindingSpecNode(packet.pkg.TypesInfo, stmtType) {
				switch expType := stmtType.X.(type) {
				case *ast.CallExpr:
					for i := 0; i < len(expType.Args); i++ {
						arg := expType.Args[i]
						if call, ok := arg.(*ast.CallExpr); ok {
							if isBindingSpecNode(packet.pkg.TypesInfo, call) {
								if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
									if id, ok := sel.X.(*ast.Ident); ok {
										if id.Name == "zipline" {
											binding := z.processStatement(packet.pkg, stmtType)

											handler := astcopy.Expr(expType.Args[i])

											expType.Args[i] = newCallExpression(binding, handler)
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
}

func (z *Zipline) processStatement(pkg *packages.Package, stmt *ast.ExprStmt) *binding {
	binding := parseSpec(pkg, stmt)
	z.renderer.render(binding)
	return binding
}

func newCallExpression(binding *binding, arg ast.Expr) *ast.CallExpr {
	ce := &ast.CallExpr{
		Fun: &ast.Ident{
			Name: binding.id() + "HandlerFunc",
		},
		Args: arg.(*ast.CallExpr).Args,
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
		binding.handler = newFromSelectorExpr(pkg, handler)
	case *ast.Ident:
		binding.handler = newFromIdent(pkg, handler)
	}

	return binding
}

func newFromSelectorExpr(pkg *packages.Package, handler *ast.SelectorExpr) *handlerInfo {
	return newFromIdent(pkg, handler.Sel)
}

func newFromIdent(pkg *packages.Package, handler *ast.Ident) *handlerInfo {
	obj := qualifiedIdentObject(pkg.TypesInfo, handler)

	sig := obj.Type().(*types.Signature)

	var id bytes.Buffer

	id.WriteString(strings.Title(obj.Pkg().Name()))
	if sig.Recv() != nil {
		rcvr := sig.Recv().Type().String()
		idx := strings.LastIndex(rcvr, ".")

		if idx > 0 {
			rcvr = strings.Trim(rcvr[idx:len(rcvr)], ".")
		}
		id.WriteString(rcvr)
	}
	id.WriteString(obj.Name())

	hi := &handlerInfo{id: string(id.Bytes())}

	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		hi.params = append(hi.params, &varToken{
			cpkg:      pkg.Name,
			name:      p.Name(),
			signature: p.Type().String(),
		})
	}

	for i := 0; i < sig.Results().Len(); i++ {
		r := sig.Results().At(i)
		hi.returns = append(hi.returns, &varToken{
			cpkg:      pkg.Name,
			name:      r.Name(),
			signature: r.Type().String(),
		})
	}

	return hi
}
