package internal

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/util"

	"github.com/bilal-bhatti/zipline/internal/tokens"
	"golang.org/x/tools/go/packages"
)

type details struct {
	lines []string
}

func newDetails(msg string) *details {
	return &details{
		lines: []string{msg},
	}
}

func (d *details) print(msg string, vals ...interface{}) {
	if len(vals) > 0 {
		d.lines = append(d.lines, fmt.Sprintf(msg, vals...))
	} else {
		d.lines = append(d.lines, msg)
	}
}

type ziplineError struct {
	msg  *details
	hint string
}

func (e ziplineError) Error() string {
	return fmt.Sprintf("%s\n\thint: %s", strings.Join(e.msg.lines, "\n\t"), e.hint)
}

func newErrorForStmt(msg string, stmt ast.Stmt) error {
	fset := token.NewFileSet()
	buf := util.NewBuffer()
	printer.Fprint(buf, fset, stmt)

	lines := newDetails(msg)

	return ziplineError{
		msg:  lines,
		hint: buf.String(),
	}
}

func newErrorForSliceVar(msg string, obj types.Object) error {
	fset := token.NewFileSet()
	buf := util.NewBuffer()
	printer.Fprint(buf, fset, obj.String())

	lines := newDetails(msg)

	return ziplineError{
		msg:  lines,
		hint: buf.String(),
	}
}

func newHandlerCallError(msg string, pkg *packages.Package, b *binding, assn *ast.AssignStmt) error {
	lines := newDetails(msg)
	lines.print("returns mismatch, expected %d, encountered %d", len(b.handler.returns), len(assn.Lhs))
	lines.print("with binding: %s", b.spec)
	return ziplineError{
		msg:  lines,
		hint: "check template and ensure app handler returns are as expected",
	}
}

func newHandlerNotResolvedError(msg string, b *binding, rets []string) error {
	lines := newDetails(msg)

	return ziplineError{
		msg:  lines,
		hint: fmt.Sprintf("ensure you have a provider function that returns %d values", len(rets)),
	}
}

func newParameterError(msg string, pkg *packages.Package, b *binding, p *tokens.TypeToken) error {
	lines := newDetails(msg)
	lines.print("missing template for parameter (%s %s)", p.VarName(), p.DeclSignature(pkg.PkgPath))
	lines.print("for function %s: %s", b.handler.sel, b.handler.signature.String())

	if b.handler.signature.Recv() != nil {
		lines.print("of receiver %s", b.handler.signature.Recv().Type().String())
	}
	lines.print("with binding: %s", b.spec)

	return ziplineError{
		msg:  lines,
		hint: fmt.Sprintf("compare parameter list and binding template, expecting %d but encountered %d", len(b.handler.params), len(b.paramTemplates)),
	}
}

func newParameterProviderError(msg string, pkg *packages.Package, b *binding, p *tokens.TypeToken) error {
	lines := newDetails(msg)
	lines.print("no provider function for parameter (%s %s)", p.VarName(), p.DeclSignature(pkg.PkgPath))
	lines.print("for function %s: %s", b.handler.sel, b.handler.signature.String())

	if b.handler.signature.Recv() != nil {
		lines.print("of receiver %s", b.handler.signature.Recv().Type().String())
	}
	lines.print("with binding: %s", b.spec)

	return ziplineError{
		msg:  lines,
		hint: fmt.Sprintf("create a provider function with a single return and \n\trequiring only parameters declared in function template body"),
	}
}

func newParameterTemplateError(msg string, pkg *packages.Package, t *template, b *binding, p *tokens.TypeToken) error {
	lines := newDetails(msg)
	lines.print("template %s doesn't support type %s", t.funcDecl.Name, p.FullSignature)
	lines.print("for function %s: %s", b.handler.sel, b.handler.signature.String())

	if b.handler.signature.Recv() != nil {
		lines.print("of receiver %s", b.handler.signature.Recv().Type().String())
	}
	lines.print("with binding: %s", b.spec)

	return ziplineError{
		msg:  lines,
		hint: fmt.Sprintf("update %s template to handle type %s", t.funcDecl.Name.String(), p.FullSignature),
	}
}
