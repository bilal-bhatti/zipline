package internal

import (
	"go/ast"
	"go/types"

	"github.com/bilal-bhatti/zipline/internal/docparser"
	"github.com/bilal-bhatti/zipline/internal/tokens"
	"golang.org/x/tools/go/packages"
)

type (
	packet struct {
		pkg      *packages.Package
		funcDecl *ast.FuncDecl
		bindings []*binding
	}

	binding struct {
		spec                   string
		method, template, path string
		handler                *handlerInfo
		paramTemplates         []string
		boundParams            []*tokens.TypeToken
	}

	handlerInfo struct {
		signature *types.Signature
		docs      *docparser.DocData
		id        string
		pkg       string
		x         *tokens.TypeToken
		sel       string
		params    []*tokens.TypeToken
		returns   []*tokens.TypeToken
	}
)

func (b binding) id() string {
	return b.handler.id
}
