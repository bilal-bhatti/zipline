package internal

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

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
		boundParams    []*typeToken
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

func (b binding) boundParamsList() string {
	params := []string{}
	for _, p := range b.boundParams {
		params = append(params, fmt.Sprintf("%s %s", p.varName(), p.param()))
	}

	return strings.Join(params, ",")
}
