package internal

import (
	"bytes"
	"strings"

	"github.com/iancoleman/strcase"
)

type varToken struct {
	cpkg, name, signature string
	isPtr                 bool
}

func newVarToken(cpkg, signature, name string) *varToken {
	vt := &varToken{
		cpkg: cpkg,
		name: name,
	}

	vt.isPtr = strings.HasPrefix(signature, "*")
	vt.signature = strings.TrimPrefix(signature, "*")

	return vt
}

func (vt varToken) sameType(signature string) bool {
	return vt.signature == strings.TrimPrefix(signature, "*")
}

func (vt varToken) varName() string {
	if vt.name != "" {
		return vt.name
	}

	vn := vt.signature
	idx := strings.LastIndex(vn, ".")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], ".")
	}

	return strcase.ToLowerCamel(vn)
}

func (vt varToken) varNameAsPointer() string {
	if vt.isPtr {
		return vt.varName()
	}

	return "&" + vt.varName()
}

func (vt varToken) inst() string {
	var b bytes.Buffer
	if vt.isPtr {
		b.WriteString("&")
	}

	vn := vt.signature
	idx := strings.LastIndex(vn, "/")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], "/")
	}

	// remove package prefix if same
	vn = strings.TrimPrefix(vn, vt.cpkg)
	vn = strings.TrimPrefix(vn, ".")

	b.WriteString(vn)
	b.WriteString("{}")

	return b.String()
}

func (vt varToken) pkg() string {
	pn := vt.signature

	idx := strings.LastIndex(pn, ".")

	if idx > 0 {
		pn = strings.Trim(pn[:idx], ".")
	}

	if pn == strings.TrimPrefix(vt.signature, "*") {
		return ""
	}

	return pn
}

func (vt varToken) param() string {
	var b bytes.Buffer
	if vt.isPtr {
		b.WriteString("*")
	}

	vn := vt.signature
	idx := strings.LastIndex(vn, "/")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], "/")
	}

	// remove package prefix if same
	vn = strings.TrimPrefix(vn, vt.cpkg)
	vn = strings.TrimPrefix(vn, ".")

	b.WriteString(vn)

	return b.String()
}
