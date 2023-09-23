package internal

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/bilal-bhatti/zipline/internal/debug"
	"github.com/bilal-bhatti/zipline/internal/util"

	"github.com/bilal-bhatti/zipline/internal/tokens"
	"github.com/go-toolsmith/astcopy"
	"golang.org/x/tools/go/packages"
)

type renderer struct {
	templates map[string]*template
	provider  *provider
	imps      []string
	preamble  *util.Buffer
	body      *util.Buffer
}

func (r *renderer) imp(imp string) {
	if imp != "" {
		r.imps = append(r.imps, "\""+imp+"\"")
	}
}

func newRenderer(templates map[string]*template, provider *provider) *renderer {
	r := &renderer{
		templates: templates,
		provider:  provider,
		imps:      make([]string, 0),
		preamble:  util.NewBuffer(),
		body:      util.NewBuffer(),
	}

	return r
}

func (r *renderer) render(pkg *packages.Package, info *binding) error {
	for _, bp := range info.boundParams {
		r.provider.memorizeTypeToken(bp)
	}

	template := r.templates[info.template]
	if template != nil {
		bites, err := r.renderFunctionTemplate(pkg, template, info)
		if err != nil {
			return err
		}
		r.body.WriteBuffer(bites)
	} else {
		return errors.Errorf("template not found for %s", info.template)
	}
	return nil
}

func pkgName(packet *packet) string {
	pkn := packet.pkg.Name
	idx := strings.LastIndex(pkn, "/")
	if idx > 0 {
		pkn = pkn[:]
		return pkn
	}
	return pkn
}

func (r *renderer) complete(packet *packet) {
	r.preamble.Sprintf("// Code generated by Zipline. DO NOT EDIT.\n\n")
	// r.preamble.Sprintf("//go:generate zipline\n")
	r.preamble.Sprintf("//+build !ziplinegen\n\n")

	// write package name
	r.preamble.Sprintf("package %s\n\n", pkgName(packet))

	// write imports
	r.preamble.Sprintf("import (\n%s)\n\n", strings.Join(r.imps, "\n"))

	// write binding func
	buf := util.NewBuffer()
	printer.Fprint(buf, packet.pkg.Fset, packet.funcDecl)
	r.preamble.Sprintf("\n\n")
	r.preamble.WriteBuffer(buf)
	r.preamble.Sprintf("\n\n")

	// write generated handler funcs
	r.preamble.WriteBuffer(r.body)
}

func (r *renderer) print(w io.Writer, frmt bool) error {
	if frmt {
		formatted, err := format.Source(r.preamble.Bytes())
		if err != nil {
			return err
		}
		w.Write(formatted)
	} else {
		w.Write(r.preamble.Bytes())
	}

	return nil
}

func (r *renderer) renderFunctionTemplate(pkg *packages.Package, t *template, b *binding) (*util.Buffer, error) {
	buf := util.NewBuffer()

	buf.Sprintf("// %s%s handles requests to:\n", b.id(), t.funcSuffix())
	buf.Sprintf("// path  : %s\n", b.path)
	buf.Sprintf("// method: %s\n", strings.ToLower(b.method))
	for _, c := range b.handler.comments.Raw {
		buf.Sprintf("// %s\n", c)
	}

	args := tokens.Join(b.boundParams, func(t *tokens.TypeToken) string {
		return t.ArgDeclaration(pkg.PkgPath)
	})

	// print function signature
	buf.Sprintf("func %s%s(%s) %s {\n", b.id(), t.funcSuffix(), args, t.returnType())

	for _, stmt := range t.funcDecl.Body.List {
		switch st := stmt.(type) {
		case *ast.ReturnStmt:
			// return statement must return a func literal
			err := r.renderFuncLiteral(pkg, b, st, buf)
			if err != nil {
				return nil, err
			}
		default:
			printer.Fprint(buf, pkg.Fset, st)
			buf.Sprintf("\n")
		}
	}

	// close func body
	buf.Sprintf("}\n\n")

	return buf, nil
}

func (r *renderer) renderFuncLiteral(pkg *packages.Package, b *binding, ret *ast.ReturnStmt, buf *util.Buffer) error {
	// statement must be a return statement
	// that returns a function literal of type
	// matching the closure return type

	// TODO: add a check and fail gracefully
	funclit := ret.Results[0].(*ast.FuncLit)

	// parse funclit type to extract type tokens
	for _, param := range funclit.Type.Params.List {
		for _, pn := range param.Names {
			r.provider.memorize(pn)
		}
	}

	// print return literal func statement
	buf.Sprintf("return ")
	printer.Fprint(buf, pkg.Fset, funclit.Type)
	buf.Sprintf(" {\n")

	// process body of template func literal
	for _, fstmt := range funclit.Body.List {
		if r.devNull(fstmt) {
			// omit
			continue
		}

		if assignStmt, ok := findZiplineNodeAs(pkg.TypesInfo, fstmt, reflect.TypeOf(&ast.AssignStmt{})); ok {
			if selectorExpr, ok := findZiplineNodeAs(pkg.TypesInfo, assignStmt, reflect.TypeOf(&ast.SelectorExpr{})); ok {
				// type coerce to make sure the find is working properly
				assign := assignStmt.(*ast.AssignStmt)
				selector := selectorExpr.(*ast.SelectorExpr)

				if selector.Sel.String() == ZiplineTemplateResolve {
					if err := r.resolve(pkg, b, assign, buf); err != nil {
						return err
					}
				} else {
					if err := r.generate(pkg, b, assign, buf); err != nil {
						return err
					}
				}

				continue
			}
		} else {
			printer.Fprint(buf, pkg.Fset, fstmt)
			buf.Sprintf("\n")
		}
	}

	buf.Sprintf("}\n")

	return nil
}

func (r *renderer) resolve(pkg *packages.Package, b *binding, assn *ast.AssignStmt, buf *util.Buffer) error {
	// resolve/print handler
	if b.handler.x != nil {
		rets := []string{}
		// extract return names from the marker call
		for _, lhs := range assn.Lhs {
			rets = append(rets, lhs.(*ast.Ident).String())
		}

		xp, err := r.provider.provide(b.handler.x, rets)
		if err != nil {
			return newHandlerNotResolvedError(err.Error(), b, rets)
		}
		buf.Sprintf("\n// initialize application handler\n")
		buf.Sprintf(xp.Call(pkg.PkgPath, assn.Tok) + "\n")
		r.imp(xp.Pkg())
	}

	return nil
}

func (r *renderer) generate(pkg *packages.Package, b *binding, assn *ast.AssignStmt, buf *util.Buffer) error {
	// resolve/print app handler dependencies and retain their var names
	params, err := r.parameters(pkg, b, buf)
	if err != nil {
		return err
	}

	if len(assn.Lhs) != len(b.handler.returns) {
		return newHandlerCallError("application handler returns mismatch", pkg, b, assn)
	}

	rets := make([]*tokens.TypeToken, len(assn.Lhs))

	// extract return names from the marker call and preserve names
	for idx := range assn.Lhs {
		lhs := assn.Lhs[idx]
		ret := b.handler.returns[idx]

		rets[idx] = tokens.NewTypeToken(ret.FullSignature, lhs.(*ast.Ident).String())
	}

	buf.Sprintf("\n// execute application handler\n")

	// invoked app handler function
	appFunc := tokens.FuncToken{
		Rets: rets,
		Args: params,
	}

	if b.handler.x != nil {
		// a struct method is the handler
		svcStruct, ok := r.provider.typeTokenFor(b.handler.x)
		if !ok {
			return errors.New("dependencies not satisfied")
		}

		// signature == handler.Something
		appFunc.Signature = fmt.Sprintf("%s.%s", svcStruct.VarName(), b.handler.sel)
	} else {
		// resolve package name for plain function handler
		// signature == packagepath.Something
		appFunc.Signature = fmt.Sprintf("%s.%s", b.handler.pkg, b.handler.sel)
	}

	// call handler with params
	call := appFunc.Call(pkg.PkgPath, assn.Tok)
	debug.Trace("generated call `%s`", call)
	buf.Sprintf("%s\n", call)

	return nil
}

func (r *renderer) parameters(pkg *packages.Package, b *binding, buf *util.Buffer) ([]*tokens.TypeToken, error) {
	params := []*tokens.TypeToken{}

	for i := 0; i < len(b.handler.params); i++ {
		p := b.handler.params[i]

		if len(b.paramTemplates) > i {
			tn := b.paramTemplates[i]
			template := r.templates[tn]
			if template != nil {
				debug.Trace("using template %s for `%s %s`\n", tn, p.VarName(), p.FullSignature)
				buf.Sprintf("\n// resolve parameter [%s] with [%s] template\n", p.VarName(), tn)
				switch tn {
				case "Query", "Path":
					err := r.renderParamTemplate(pkg, template, b, p, buf)
					if err != nil {
						return nil, err
					}
				case "Body":
					err := r.renderBodyTemplate(pkg, template, b, p, buf)
					if err != nil {
						return nil, err
					}
				default:
					err := r.renderGenericTemplate(pkg, template, b, p, buf)
					if err != nil {
						return nil, err
					}
				}

				params = append(params, p)
				buf.Sprintf("\n")
				continue
			}

			if tn == ZiplineTemplateResolve {
				debug.Trace("find a provider function for `%s %s`\n", p.VarName(), p.FullSignature)
				ft, err := r.provider.provide(p, []string{p.VarName()})
				if err != nil {
					// check if type is already known
					if known, ok := r.provider.typeTokenFor(p); ok {
						debug.Trace("known variable `%s %s` as `%s`\n", p.VarName(), p.FullSignature, known.VarName())
						// use var type token declared in the template body
						params = append(params, known)
						continue
					}
					return nil, newParameterProviderError("failed to resolve handler parameters", pkg, b, p)
				}

				buf.Sprintf("\n// resolve parameter [%s] through a provider\n", p.VarName())
				buf.Sprintf(ft.Call(pkg.PkgPath, token.DEFINE) + "\n")

				params = append(params, p)
				continue
			}
		}

		// faild to find a way to satisfy parameter
		return nil, newParameterError("failed to resolve handler parameters", pkg, b, p)
	}

	return params, nil
}

func (r *renderer) renderParamTemplate(pkg *packages.Package, t *template, b *binding, p *tokens.TypeToken, buf *util.Buffer) error {
	var tmplBody []ast.Stmt

	if len(t.funcDecl.Body.List) != 1 {
		return errors.New(fmt.Sprintf("template %s must contain a single switch statement", t.funcDecl.Name))
	}

	switchStmt, ok := t.funcDecl.Body.List[0].(*ast.SwitchStmt)
	if !ok {
		return errors.New(fmt.Sprintf("template %s must contain a single switch statement", t.funcDecl.Name))
	}

	for _, sb := range switchStmt.Body.List {
		caseStmt, ok := sb.(*ast.CaseClause)
		if !ok {
			return errors.New(fmt.Sprintf("template %s must contain a signle switch statement", t.funcDecl.Name))
		}

		for _, csl := range caseStmt.List {
			blit, ok := csl.(*ast.BasicLit)
			if !ok {
				continue
			}

			if blit.Value == fmt.Sprintf("\"%s\"", p.FullSignature) {
				tmplBody = caseStmt.Body
				break
			}
		}
	}

	if tmplBody == nil {
		return newParameterTemplateError("failed to locate request parameters", pkg, t, b, p)
	}

	var format string
	tags, ok := b.handler.comments.Tags[p.VarName()]
	if ok {
		tag, err := tags.Get("format")
		if err == nil {
			format = fmt.Sprintf("%s,%s", tag.Name, strings.Join(tag.Options, ","))
		}
	}

	tbuf := r.rename("name", format, pkg, b, p, tmplBody)
	buf.WriteBuffer(tbuf)

	return nil
}

func (r *renderer) renderBodyTemplate(pkg *packages.Package, t *template, b *binding, p *tokens.TypeToken, buf *util.Buffer) error {
	return r.renderGenericTemplate(pkg, t, b, p, buf)
}

func (r *renderer) renderGenericTemplate(pkg *packages.Package, t *template, b *binding, p *tokens.TypeToken, buf *util.Buffer) error {
	tbuf := r.rename("name", "", pkg, b, p, t.funcDecl.Body.List)
	buf.WriteBuffer(tbuf)
	return nil
}

func (r *renderer) rename(old, format string, pkg *packages.Package, b *binding, new *tokens.TypeToken, body []ast.Stmt) *util.Buffer {
	renamer := func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.BasicLit:
			if n.Value == fmt.Sprintf("\"%s\"", old) {
				n.Value = fmt.Sprintf("\"%s\"", new.Name)
			}
			if format != "" && n.Value == fmt.Sprintf("\"format\"") {
				n.Value = fmt.Sprintf("\"%s\"", format)
			}
		case *ast.Ident:
			if n.Name == old {
				n.Name = new.Name
			}
		case *ast.CallExpr:
			for i := 0; i < len(n.Args); i++ {
				arg := n.Args[i]
				ue, ok := arg.(*ast.UnaryExpr)
				if !ok {
					return true
				}

				id, ok := ue.X.(*ast.Ident)
				if !ok {
					return true
				}

				if id.Name == old {
					if new.IsPtr {
						// param was declared as a pointer
						// replace unary expression with ident
						n.Args[i] = ast.NewIdent(new.Name)
					}
				}
			}
		}

		return true
	}

	buf := util.NewBuffer()

	for _, stmt := range body {
		stmtCopy := astcopy.Stmt(stmt)

		switch stmt := stmt.(type) {
		case *ast.DeclStmt:
			// omit printing var decls that are repeated
			gd, ok := stmt.Decl.(*ast.GenDecl)
			if !ok {
				goto Include
			}

			// only if a single var declaration
			if gd.Tok != token.VAR || len(gd.Specs) != 1 {
				goto Include
			}

			vs, ok := gd.Specs[0].(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 {
				goto Include
			}

			if id, ok := vs.Type.(*ast.Ident); ok {
				// mostly to deal with duplidate
				// var err error declarations
				// id.Name ~=~ error
				v, ok := r.provider.known[id.Name]
				// omit if known var with same name
				if ok && v.Name == vs.Names[0].String() {
					debug.Trace("omitting `var %s %s` from generated code, already declared", v.Name, id.String())
					continue
				}
			}
		case *ast.AssignStmt:
			if r.newStructValue(stmt) {
				// get new object
				inst := fmt.Sprintf("%s %s %s", new.VarName(), stmt.Tok.String(), new.NewInstance(pkg.PkgPath))
				debug.Trace("new struct `%s`", inst)
				buf.Sprintf("%s\n", inst)
				continue
			}
			ast.Inspect(stmtCopy, renamer)
		case *ast.ExprStmt:
			// filter out zipline directive
			if r.devNull(stmt) {
				continue
			}
			ast.Inspect(stmtCopy, renamer)
		default:
			ast.Inspect(stmtCopy, renamer)
		}
	Include:
		printer.Fprint(buf, token.NewFileSet(), stmtCopy)
		buf.Sprintf("\n")
	}

	return buf
}

func (r *renderer) devNull(stmt ast.Stmt) bool {
	expr, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return false
	}

	call, ok := expr.X.(*ast.CallExpr)
	if !ok {
		return false
	}

	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		zobj := r.provider.qualifiedIdentObject(sel.X)
		if zobj != nil && strings.HasSuffix(zobj.Type().String(), ZiplineTemplate) {
			// if ignore
			if sel.Sel.Name == ZiplineTemplateDevNull || sel.Sel.Name == ZiplineTemplateIgnore {
				return true
			}
		}
	}

	return false
}

func (r *renderer) newStructValue(expr *ast.AssignStmt) bool {
	// check if it's a new ZiplineTemplate object
	// must be a simple assignment like
	// name := ZiplineTemplate{} or name := &ZiplineTemplate{}
	if len(expr.Lhs) != 1 && len(expr.Rhs) != 1 {
		return false
	}

	var newObj *ast.CompositeLit

	if unary, ok := expr.Rhs[0].(*ast.UnaryExpr); ok {
		newObj, ok = unary.X.(*ast.CompositeLit)
		if !ok {
			return false
		}
	} else {
		newObj, ok = expr.Rhs[0].(*ast.CompositeLit)
		if !ok {
			return false
		}
	}

	valueType, ok := newObj.Type.(*ast.Ident)
	if !ok || valueType.Name != ZiplineTemplate {
		return false
	}

	return true
}
