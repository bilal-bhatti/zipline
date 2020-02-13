package internal

import (
	"go/ast"
)

type template struct {
	funcDecl *ast.FuncDecl
}

func newTemplate(f *ast.FuncDecl) *template {
	return &template{
		funcDecl: f,
	}
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
