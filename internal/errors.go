package internal

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
)

type ziplineError struct {
	msg, hint string
}

func (e ziplineError) Error() string {
	return fmt.Sprintf("%s:\n\t%s", e.msg, e.hint)
}

func newErrorForStmt(msg string, stmt ast.Stmt) error {
	fset := token.NewFileSet()
	buf := newBuffer()
	printer.Fprint(buf.buf, fset, stmt)

	return ziplineError{
		msg:  msg,
		hint: buf.buf.String(),
	}
}

func newErrorForSliceVar(msg string, obj types.Object) error {
	fset := token.NewFileSet()
	buf := newBuffer()
	printer.Fprint(buf.buf, fset, obj.String())

	return ziplineError{
		msg:  msg,
		hint: buf.buf.String(),
	}
}
