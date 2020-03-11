package internal

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type (
	packet struct {
		pkg      *packages.Package
		funcDecl *ast.FuncDecl
		bindings []*binding
	}

	binding struct {
		spec           string
		template, path string
		handler        *handlerInfo
		paramTemplates []string
	}

	handlerInfo struct {
		signature *types.Signature
		comments  *comments
		id        string
		pkg       string
		x         *typeToken
		sel       string
		params    []*typeToken
		returns   []*typeToken
	}
)

func (b binding) id() string {
	return b.handler.id
}
