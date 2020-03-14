package internal

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

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
		spec           string
		template, path string
		handler        *handlerInfo
		paramTemplates []string
		boundParams    []*tokens.TypeToken
	}

	handlerInfo struct {
		signature *types.Signature
		comments  *comments
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

func (b binding) boundParamsList(importingPkg string) string {
	params := []string{}
	for _, p := range b.boundParams {
		params = append(params, fmt.Sprintf("%s %s", p.VarName(), p.DeclSignature(importingPkg)))
	}

	return strings.Join(params, ",")
}
