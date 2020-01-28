package internal

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/packages"
)

type (
	binding struct {
		name, method, path string
		handler            *handlerInfo
	}

	handlerInfo struct {
		params  []*BindingVar
		returns []*BindingVar
	}

	BindingVar struct {
		Name, Type string
	}

	Call struct {
		Package, Receiver, Message string
		IsPointer                  bool
	}

	Binding struct {
		Method, Path string
		Params       []*BindingVar
		Returns      []*BindingVar
	}

	HandlerSignature struct {
		Name    string
		Params  map[string]string
		Returns map[string]string
	}
	Handler struct {
		Path             string
		HttpMethod       string
		HandlerSignature *HandlerSignature
	}
	GenMetadata struct {
		PackageName string
		Handler     *Handler
	}
	StdGenPacket struct {
		Pkg         *packages.Package
		Bindings    *ast.FuncDecl
		GenMetadata []*GenMetadata
		Preamble    string
		Specs       []*ast.ExprStmt
	}

	StdPackageParser struct {
		Packets []*StdGenPacket
	}
)

func NewBinding(m, p string) *Binding {
	return &Binding{
		Method:  m,
		Path:    p,
		Params:  []*BindingVar{},
		Returns: []*BindingVar{},
	}
}

func (b *Binding) AddParam(n, t string) {
	b.Params = append(b.Params, &BindingVar{
		Name: n,
		Type: t,
	})
}

func (b *Binding) AddReturn(n, t string) {
	b.Returns = append(b.Returns, &BindingVar{
		Name: n,
		Type: t,
	})
}

func (bv BindingVar) typeInfo() (bool, string, string) {
	split := strings.Split(bv.Type, ".")

	pointer := false
	if strings.HasPrefix(split[0], "*") {
		pointer = true
		split[0] = strings.TrimPrefix(split[0], "*")
	}

	if len(split) == 0 {
		return pointer, "", ""
	}
	if len(split) == 1 {
		return pointer, split[0], ""
	}
	return pointer, split[0], split[1]
}

func (b binding) id() string {
	return "XYZ"
}
