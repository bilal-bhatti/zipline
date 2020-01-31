package internal

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"strings"
)

type renderer struct {
	imports  []string
	types    []string
	preamble *bytes.Buffer
	body     *bytes.Buffer
}

func newRenderer() *renderer {
	r := &renderer{
		imports:  make([]string, 0),
		types:    make([]string, 0),
		preamble: &bytes.Buffer{},
		body:     &bytes.Buffer{},
	}

	// r.ws("package services\n\n")
	return r
}

func (r *renderer) render(info *binding) {
	switch info.method {
	case "Post":
		r.post(info)
	case "Get":
	}
}

func (r *renderer) ws(s string, vals ...interface{}) {
	if len(vals) > 0 {
		r.body.WriteString(fmt.Sprintf(s, vals...))
	} else {
		r.body.WriteString(s)
	}
}

func (r *renderer) complete() {
	r.preamble.WriteString("package services\n\n")
	r.preamble.WriteString("import (")
	for _, imp := range r.imports {
		r.preamble.WriteString(fmt.Sprintf("\"%s\"\n", imp))
	}
	r.preamble.WriteString(")\n\n")

	r.preamble.WriteString("type (\n")
	for _, tipe := range r.types {
		r.preamble.WriteString(fmt.Sprintf("%s\n", tipe))
	}
	r.preamble.WriteString(")\n\n")

	r.preamble.Write(r.body.Bytes())
}

func (r *renderer) print(frmt bool) {
	if frmt {
		formatted, err := format.Source(r.preamble.Bytes())
		if err != nil {
			panic(err)
		}
		fmt.Print("\n", string(formatted))
	} else {
		fmt.Print("\n", string(r.preamble.Bytes()))
	}
}

func (r *renderer) post(info *binding) {
	r.recordFuncType(info)
	r.ws("func %sHandlerFunc(funk %sHandlerType) http.HandlerFunc {\n", info.id(), info.id())
	r.ws("return func(w http.ResponseWriter, req *http.Request) {\n")

	varNames := []string{}
	for _, p := range info.handler.params {
		log.Println("signature", p.signature)
		switch p.signature {
		// TODO: use context provider
		case "context.Context":
			r.ws(fmt.Sprintf("%s := context.TODO()\n", p.varName()))
			varNames = append(varNames, p.varName())
		default:
			varNames = append(varNames, r.writeJsonDecoder(p))
		}

		if p.pkg() != "" {
			r.imports = append(r.imports, p.pkg())
		}
		r.ws("\n")
	}

	r.ws("res, err := funk(%s)\n", strings.Join(varNames, ","))

	r.ws("if err != nil {\n")
	r.ws("// write error response\n")
	r.ws("// internal error\n")
	r.ws("panic(err)\n")
	r.ws("}\n\n")

	r.ws("w.WriteHeader(http.StatusOK)\n")
	r.ws("w.Header().Set(\"Content-Type\", \"text/plain; charset=utf-8\")\n")
	r.ws("err = json.NewEncoder(w).Encode(res)\n")
	r.ws("if err != nil {\n")
	r.ws("// write error response\n")
	r.ws("panic(err)\n")
	r.ws("}\n")

	r.ws("}\n")
	r.ws("}\n\n")
}

func (r *renderer) writeJsonDecoder(p *varToken) string {

	r.writeAssignment(p)

	r.ws("err := json.NewDecoder(req.Body).Decode(%s)\n", p.varNameAsPointer())

	r.ws("if err != nil {\n")
	r.ws("// write error response\n")
	r.ws("// invalid request error\n")
	r.ws("panic(err)\n")
	r.ws("}\n\n")

	return p.varName()
}

func (r *renderer) writeAssignment(p *varToken) string {
	r.ws("%s := %s\n", p.varName(), p.inst())

	return p.varName()
}

func (r *renderer) recordFuncType(b *binding) {
	buf := bytes.Buffer{}

	buf.WriteString("XYZHandlerType func(")
	for i := 0; i < len(b.handler.params); i++ {
		p := b.handler.params[i]

		buf.WriteString(fmt.Sprintf("%s %s", p.varName(), p.param()))

		if i+1 < len(b.handler.params) {
			buf.WriteString(", ")
		}
	}

	buf.WriteString(") (")
	for i := 0; i < len(b.handler.returns); i++ {
		r := b.handler.returns[i]

		buf.WriteString(r.param())

		if i+1 < len(b.handler.returns) {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(")\n")

	r.types = append(r.types, string(buf.Bytes()))
}
