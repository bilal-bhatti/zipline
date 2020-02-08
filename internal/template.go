package internal

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type template struct {
	funcDecl *ast.FuncDecl
}

func newTemplate(f *ast.FuncDecl) *template {
	return &template{
		funcDecl: f,
	}
}

func loadTemplates(pkgs []*packages.Package) map[string]*template {
	templates := make(map[string]*template)
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, d := range file.Decls {
				if f, ok := d.(*ast.FuncDecl); ok {
					if f.Recv != nil && len(f.Recv.List) == 1 {
						// match ZiplineTemplate as receiver
						// record receiver object identity
						obj := f.Recv.List[0]
						if zt, ok := obj.Type.(*ast.Ident); ok {
							if zt.String() == "ZiplineTemplate" {
								templates[f.Name.String()] = newTemplate(f)
							}
						}
					}
				}
			}
		}
	}
	return templates
}

func (t template) returnType() string {
	ret := t.funcDecl.Type.Results.List[0]
	switch rt := ret.Type.(type) {
	case *ast.SelectorExpr:
		var retType string
		if rt.X != nil {
			retType = rt.X.(*ast.Ident).String() + "."
		}
		return retType + rt.Sel.String()
	}
	return "http.HandlerFunck"
}

func (t template) funcSuffix() string {
	ret := t.funcDecl.Type.Results.List[0]
	switch rt := ret.Type.(type) {
	case *ast.SelectorExpr:
		return rt.Sel.String()
	}
	return "HandlerFunck"
}
