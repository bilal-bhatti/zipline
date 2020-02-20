package internal

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
)

type ziplineError struct {
	msg, stmt string
}

func (e ziplineError) Error() string {
	return fmt.Sprintf("%s:\n\t%s", e.msg, e.stmt)
}

func newError(msg string, stmt ast.Stmt) error {
	fset := token.NewFileSet()
	buf := newBuffer()
	printer.Fprint(buf.buf, fset, stmt)

	return ziplineError{
		msg:  msg,
		stmt: buf.buf.String(),
	}
}
