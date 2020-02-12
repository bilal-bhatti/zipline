package internal

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"io"
	"log"
	"strings"
)

type renderer struct {
	templates map[string]*template
	provider  *provider
	imps      []string
	preamble  buffer
	body      buffer
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
		preamble:  newBuffer(),
		body:      newBuffer(),
	}

	return r
}

func (r *renderer) render(info *binding) {
	template := r.templates[info.template]
	if template != nil {
		bites, err := r.renderTemplate(template, info)
		if err != nil {
			panic(err)
		}
		r.body.buf.Write(bites)
	} else {
		log.Println("Template not found for", info.template)
	}
}

func pkgName(packet *packet) string {
	pkn := packet.pkg.Name
	idx := strings.LastIndex(pkn, "/")
	if idx > 0 {
		pkn = pkn[:len(pkn)]
		return pkn
	}
	return pkn
}

func (r *renderer) complete(packet *packet) {
	r.preamble.ws("// Code generated by Zipline. DO NOT EDIT.\n\n")

	r.preamble.ws("//go:generate zipline\n")
	r.preamble.ws("// +build !ziplinegen\n\n")

	// write package name
	r.preamble.ws("package %s\n\n", pkgName(packet))

	// write imports
	r.preamble.ws("import (\n%s)\n\n", strings.Join(r.imps, "\n"))

	// write binding func
	fset := token.NewFileSet()
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, packet.funcDecl)
	r.preamble.ws("\n\n")
	r.preamble.buf.Write(buf.Bytes())
	r.preamble.ws("\n\n")

	// write generated handler funcs
	r.preamble.buf.Write(r.body.buf.Bytes())
}

func (r *renderer) print(w io.Writer, frmt bool) {
	if frmt {
		formatted, err := format.Source(r.preamble.buf.Bytes())
		if err != nil {
			panic(err)
		}
		w.Write(formatted)
	} else {
		w.Write(r.preamble.buf.Bytes())
	}
}

func (r *renderer) renderTemplate(t *template, b *binding) ([]byte, error) {
	fset := token.NewFileSet()

	buf := newBuffer()

	for _, stmt := range t.funcDecl.Body.List {
		// must be a return statement
		ret := stmt.(*ast.ReturnStmt)

		funclit := ret.Results[0].(*ast.FuncLit)
		// parse funclit type to extract type tokens
		for _, param := range funclit.Type.Params.List {
			for _, pn := range param.Names {
				obj := r.provider.qualifiedIdentObject(pn)
				tt := newTypeToken("", obj.Type().String(), obj.Name())
				r.provider.known[tt.signature] = tt
			}
		}

		buf.ws("// %s%s handles requests to:\n", b.id(), t.funcSuffix())
		buf.ws("// path  : %s\n", b.path)
		buf.ws("// method: %s\n", strings.ToLower(b.template))
		buf.ws("func %s%s() %s {\n", b.id(), t.funcSuffix(), t.returnType())
		buf.ws("return ")
		printer.Fprint(buf.buf, fset, funclit.Type)
		buf.ws(" {\n")

		for _, fstmt := range funclit.Body.List {
			var expand = false

			switch stmtType := fstmt.(type) {
			case *ast.AssignStmt:
				if call, ok := stmtType.Rhs[0].(*ast.CallExpr); ok {
					if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
						obj := r.provider.qualifiedIdentObject(selector.X)
						if obj != nil && strings.HasSuffix(obj.Type().String(), ZiplineTemplate) {
							expand = true
							r.expand(b, stmtType, buf)
						}
					}
				}
				if !expand {
					printer.Fprint(buf.buf, fset, fstmt)
					buf.ws("\n")
				}
			default:
				printer.Fprint(buf.buf, fset, fstmt)
				buf.ws("\n")
			}
		}
		buf.ws("}\n")
		buf.ws("}\n\n")
	}

	return buf.buf.Bytes(), nil
}

func (r *renderer) expand(b *binding, assn *ast.AssignStmt, buf buffer) {
	// resolve/print app handler dependencies and retain their var names
	params := r.deps(b, buf)

	call := newBuffer()

	rets := []string{}
	// extract return names from the marker call
	for _, lhs := range assn.Lhs {
		rets = append(rets, lhs.(*ast.Ident).String())
	}

	call.ws(strings.Join(rets, ","))
	call.ws(" %s ", assn.Tok.String())

	if b.handler.x != nil {
		xp := r.provider.provide(b.handler.x)
		buf.ws("// initialize application handler\n")
		buf.ws(xp.call() + "\n")
		r.imp(xp.pkg())

		funk, ok := r.provider.varFor(b.handler.x)
		if !ok {
			panic("Dependencies not satisfied")
		}

		call.ws("%s.", funk.varName())
	}
	call.ws("%s(%s)\n", b.handler.sel, strings.Join(params, ","))

	buf.ws("// execute application handler\n")
	buf.ws(call.buf.String())
}

func (r *renderer) deps(b *binding, buf buffer) []string {
	hreq, _ := r.provider.varFor(HREQ)
	hwri, _ := r.provider.varFor(HWRI)

	params := []string{}
	for _, p := range b.handler.params {

		if isPathParam(b, p) {
			buf.ws("\n// parse path parameter %s\n", p.varName())
			switch p.signature {
			case "int":
				buf.ws("%s, err := strconv.Atoi(chi.URLParam(%s, \"%s\"))\n", p.name, hreq.varName(), p.name)
				buf.ws("if err != nil {\n")
				buf.ws("  // invalid request error\n")
				buf.ws("  http.Error(%s, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)\n", hwri.varName())
				buf.ws("  return\n")
				buf.ws("}\n\n")
			case "string":
				buf.ws("%s := chi.URLParam(%s, \"%s\")\n\n", p.name, hreq.varName(), p.name)
			}
		} else if ft := r.provider.provide(p); ft != nil {
			buf.ws("\n// resolve %s dependency through a provider function\n", p.varName())

			buf.ws("%s\n\n", ft.call())
		} else if b.template == "Post" || b.template == "Put" {
			buf.ws("\n// extract json body and marshall %s\n", p.varName())

			buf.ws("%s := %s\n", p.varName(), p.inst())
			buf.ws("err = json.NewDecoder(%s.Body).Decode(%s)\n", hreq.varName(), p.varNameAsPointer())
			buf.ws("if err != nil {\n")
			buf.ws("  // invalid request error\n")
			buf.ws("  http.Error(%s, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)\n", hwri.varName())
			buf.ws("  return\n")
			buf.ws("}\n\n")
		}

		params = append(params, p.varName())

		r.imp(p.pkg())

		buf.ws("\n")
	}

	return params
}

func join(tokens []*typeToken) string {
	s := []string{}
	for _, token := range tokens {
		s = append(s, token.param())
	}

	return strings.Join(s, ",")
}
