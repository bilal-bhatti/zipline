package tokens

import (
	"bytes"
	"fmt"
	"go/types"
	"strings"

	"github.com/iancoleman/strcase"
)

type TypeToken struct {
	FullSignature  string
	Name           string
	Signature      string // temp export
	IsPtr, IsSlice bool
	VarType        *types.Var
}

func (tt TypeToken) String() string {
	return fmt.Sprintf("{name: %s, type: %s}", tt.VarName(), tt.FullSignature)
}

func NewTypeToken(signature, name string) *TypeToken {
	tt := &TypeToken{
		FullSignature: signature,
		Name:          name,
	}

	tt.IsSlice = strings.HasPrefix(signature, "[]")
	tt.Signature = strings.TrimPrefix(signature, "[]")

	tt.IsPtr = strings.HasPrefix(tt.Signature, "*")
	tt.Signature = strings.TrimPrefix(tt.Signature, "*")

	return tt
}

func (tt *TypeToken) TypeTokenAs(as *TypeToken) *TypeToken {
	if tt.FullSignature == as.FullSignature {
		return tt
	}

	return NewTypeToken(as.FullSignature, as.Name)
}

func (tt TypeToken) SameType(t types.Type, strict bool) bool {
	if strict {
		return tt.FullSignature == t.String()
	}
	return tt.FullSignature == strings.TrimPrefix(t.String(), "*")
}

func (tt TypeToken) VarName() string {
	if tt.Name != "" {
		return tt.Name
	}

	vn := tt.Signature
	idx := strings.LastIndex(vn, ".")

	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], ".")
	}

	return strcase.ToLowerCamel(vn)
}

func (tt TypeToken) NewInstance(importingPkg string) string {
	var b bytes.Buffer
	// if requested type is a pointer, let's declare it as such
	if tt.IsPtr {
		b.WriteString("&")
	}

	b.WriteString(tt.SimpleSignature(importingPkg))

	b.WriteString("{}")

	return b.String()
}

func (tt TypeToken) DeclSignature(importingPkg string) string {
	var b bytes.Buffer
	if tt.IsPtr {
		b.WriteString("*")
	}

	b.WriteString(tt.SimpleSignature(importingPkg))

	return b.String()
}

// returns declaration signature without pointer details
func (tt TypeToken) SimpleSignature(importingPkg string) string {
	var b bytes.Buffer

	vn := tt.Signature

	// remove package prefix if same package
	// github.com/bilal-bhatti/zipline/example/web.EchoRequest
	// github.com/bilal-bhatti/zipline/example/web
	// EchoRequest
	vn = strings.TrimPrefix(vn, importingPkg)
	vn = strings.TrimPrefix(vn, ".")

	// if not same package, lob off up to the last /
	// github.com/bilal-bhatti/zipline/example/models.ThingResponse
	// github.com/bilal-bhatti/zipline/example/web
	// models.ThingResponse
	idx := strings.LastIndex(vn, "/")
	if idx > 0 {
		vn = strings.Trim(vn[idx:len(vn)], "/")
	}

	b.WriteString(vn)

	return b.String()
}
