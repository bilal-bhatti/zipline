package internal

import "go/ast"

type (
	binding struct {
		name, method, path string
		handler            *handlerInfo
		zipline            *ast.CallExpr
	}

	handlerInfo struct {
		id      string
		x       *varToken
		sel     string
		params  []*varToken
		returns []*varToken
	}
)

func (b binding) id() string {
	return b.handler.id
}
