package internal

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"
)

type ziplineError struct {
	msg, hint string
}

func (e ziplineError) Error() string {
	return fmt.Sprintf("%s\n\t%s", e.msg, e.hint)
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

func newHandlerNotResolvedError(msg string, b *binding, rets []string) error {
	lines := []string{msg}

	return ziplineError{
		msg:  strings.Join(lines, "\n\t"),
		hint: fmt.Sprintf("hint: ensure you have a provider function that returns %d value", len(rets)),
	}
}

func newParameterError(msg string, b *binding, p *typeToken) error {
	lines := []string{msg}
	lines = append(lines, fmt.Sprintf("missing template for parameter (%s %s)", p.varName(), p.signature))
	lines = append(lines, fmt.Sprintf("for function %s: %s", b.handler.sel, b.handler.signature.String()))

	if b.handler.signature.Recv() != nil {
		lines = append(lines, fmt.Sprintf("of receiver %s", b.handler.signature.Recv().Type().String()))
	}
	lines = append(lines, fmt.Sprintf("with binding: %s", b.spec))

	return ziplineError{
		msg:  strings.Join(lines, "\n\t"),
		hint: fmt.Sprintf("hint: compare function parameter list and binding template, expecting %d but encountered %d", len(b.handler.params), len(b.paramTemplates)),
	}
}

func newParameterProviderError(msg string, b *binding, p *typeToken) error {
	lines := []string{msg}
	lines = append(lines, fmt.Sprintf("no provider function for parameter (%s %s)", p.varName(), p.signature))
	lines = append(lines, fmt.Sprintf("for function %s: %s", b.handler.sel, b.handler.signature.String()))

	if b.handler.signature.Recv() != nil {
		lines = append(lines, fmt.Sprintf("of receiver %s", b.handler.signature.Recv().Type().String()))
	}
	lines = append(lines, fmt.Sprintf("with binding: %s", b.spec))

	return ziplineError{
		msg:  strings.Join(lines, "\n\t"),
		hint: fmt.Sprintf("hint: create a provider function with a single return and \n\trequiring only parameters declared in function template body"),
	}
}

func newParameterTemplateError(msg string, t *template, b *binding, p *typeToken) error {
	lines := []string{msg}
	lines = append(lines, fmt.Sprintf("template %s doesn't support type %s", t.funcDecl.Name, p.fullSignature()))
	lines = append(lines, fmt.Sprintf("for function %s: %s", b.handler.sel, b.handler.signature.String()))

	if b.handler.signature.Recv() != nil {
		lines = append(lines, fmt.Sprintf("of receiver %s", b.handler.signature.Recv().Type().String()))
	}
	lines = append(lines, fmt.Sprintf("with binding: %s", b.spec))

	return ziplineError{
		msg:  strings.Join(lines, "\n\t"),
		hint: fmt.Sprintf("hint: update %s template to handle type %s", t.funcDecl.Name.String(), p.fullSignature()),
	}
}
