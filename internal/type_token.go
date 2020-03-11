package internal

import (
	"bytes"
	"fmt"
	"go/types"
	"strings"

	"github.com/iancoleman/strcase"
)

type typeToken struct {
	cpkg, name, signature string
	isPtr, isSlice        bool
	varType               *types.Var
}

func (tt typeToken) String() string {
	return fmt.Sprintf("{name: %s, type: %s, slice: %v, ptr: %v}", tt.varName(), tt.signature, tt.isSlice, tt.isPtr)
}

func newTypeToken(cpkg, signature, name string) *typeToken {
	tt := &typeToken{
		cpkg: cpkg,
		name: name,
	}

	tt.isSlice = strings.HasPrefix(signature, "[]")
	tt.signature = strings.TrimPrefix(signature, "[]")

	tt.isPtr = strings.HasPrefix(tt.signature, "*")
	tt.signature = strings.TrimPrefix(tt.signature, "*")

	return tt
}

func (tt typeToken) sameType(signature string) bool {
	return tt.signature == strings.TrimPrefix(signature, "*")
}

func (tt typeToken) varName() string {
	if tt.name != "" {
		return tt.name
	}

	vn := tt.signature
	idx := strings.LastIndex(vn, ".")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], ".")
	}

	return strcase.ToLowerCamel(vn)
}

func (tt typeToken) varNameAsPointer() string {
	if tt.isPtr {
		return tt.varName()
	}

	return "&" + tt.varName()
}

func (tt typeToken) inst() string {
	var b bytes.Buffer
	if tt.isPtr {
		b.WriteString("&")
	}

	vn := tt.signature
	idx := strings.LastIndex(vn, "/")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], "/")
	}

	// remove package prefix if same
	vn = strings.TrimPrefix(vn, tt.cpkg)
	vn = strings.TrimPrefix(vn, ".")

	b.WriteString(vn)
	b.WriteString("{}")

	return b.String()
}

func (tt typeToken) pkg() string {
	pn := tt.signature

	idx := strings.LastIndex(pn, ".")

	if idx > 0 {
		pn = strings.Trim(pn[:idx], ".")
	}

	if pn == strings.TrimPrefix(tt.signature, "*") {
		return ""
	}

	return pn
}

func (tt typeToken) param() string {
	var b bytes.Buffer
	if tt.isPtr {
		b.WriteString("*")
	}

	vn := tt.signature
	idx := strings.LastIndex(vn, "/")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], "/")
	}

	// remove package prefix if same
	vn = strings.TrimPrefix(vn, tt.cpkg)
	vn = strings.TrimPrefix(vn, ".")

	b.WriteString(vn)

	return b.String()
}

func (tt typeToken) shortSignature() string {
	var b bytes.Buffer

	vn := tt.signature
	idx := strings.LastIndex(vn, "/")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], "/")
	}

	// remove package prefix if same
	vn = strings.TrimPrefix(vn, ".")

	b.WriteString(vn)

	return b.String()
}

func (tt typeToken) fullSignature() string {
	var b bytes.Buffer

	if tt.isSlice {
		b.WriteString("[]")
	}
	if tt.isPtr {
		b.WriteString("*")
	}
	b.WriteString(tt.signature)

	return b.String()
}
