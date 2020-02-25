package internal

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type (
	packet struct {
		pkg      *packages.Package
		funcDecl *ast.FuncDecl
		bindings []*binding
	}

	binding struct {
		template, path string
		handler        *handlerInfo
		paramTemplates []string
	}

	handlerInfo struct {
		comments []string
		id       string
		pkg      string
		x        *typeToken
		sel      string
		params   []*typeToken
		returns  []*typeToken
	}
)

func (b binding) id() string {
	return b.handler.id
}
