package internal

import (
	"bytes"
	"strings"

	"github.com/iancoleman/strcase"
)

type varToken struct {
	name, signature string
}

func (vt varToken) isPointer() bool {
	return strings.HasPrefix(vt.signature, "*")
}

func (vt varToken) varName() string {
	if vt.name != "" {
		return vt.name
	}

	vn := strings.TrimPrefix(vt.signature, "*")
	idx := strings.LastIndex(vn, ".")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], ".")
	}

	return strcase.ToLowerCamel(vn)
}

func (vt varToken) varNameAsPointer() string {
	if vt.isPointer() {
		return vt.varName()
	}

	return "&" + vt.varName()
}

func (vt varToken) inst() string {
	var b bytes.Buffer
	if vt.isPointer() {
		b.WriteString("&")
	}

	vn := strings.TrimPrefix(vt.signature, "*")
	idx := strings.LastIndex(vn, "/")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], "/")
	}

	b.WriteString(vn)
	b.WriteString("{}")

	return b.String()
}

func (vt varToken) pkg() string {
	pn := strings.TrimPrefix(vt.signature, "*")

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
	if vt.isPointer() {
		b.WriteString("*")
	}

	vn := strings.TrimPrefix(vt.signature, "*")
	idx := strings.LastIndex(vn, "/")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], "/")
	}

	b.WriteString(vn)

	return b.String()
}
