package tokens

import (
	"bytes"
	"fmt"
	"strings"
)

type FuncToken struct {
	Signature string
	Args      []*TypeToken
	Rets      []*TypeToken
}

func (ft FuncToken) Pkg() string {
	idx := strings.LastIndex(ft.Signature, ".")

	var pn string
	if idx > 0 {
		pn = strings.Trim(ft.Signature[:idx], ".")
	}

	if pn == ft.Signature {
		return ""
	}

	return pn
}

func (ft FuncToken) Call(importingPkg string) string {
	buf := &bytes.Buffer{}

	var fn string

	// remove package prefix if same
	fn = strings.TrimPrefix(ft.Signature, importingPkg)
	fn = strings.TrimPrefix(fn, ".")

	// if above didn't do the job
	idx := strings.LastIndex(fn, "/")

	if idx > 0 {
		fn = strings.Trim(ft.Signature[idx:len(ft.Signature)], "/")
	}

	args := []string{}
	for _, arg := range ft.Args {
		args = append(args, arg.VarName())
	}

	rets := []string{}
	for _, ret := range ft.Rets {
		rets = append(rets, ret.VarName())
	}

	buf.WriteString(fmt.Sprintf("%s := %s(%s)", strings.Join(rets, ","), fn, strings.Join(args, ",")))

	return buf.String()
}
