package internal

import (
	"strings"
)

type funcToken struct {
	cpkg, signature string
	args            []*typeToken
	rets            []*typeToken
}

func (ft funcToken) pkg() string {
	idx := strings.LastIndex(ft.signature, ".")

	var pn string
	if idx > 0 {
		pn = strings.Trim(ft.signature[:idx], ".")
	}

	if pn == ft.signature {
		return ""
	}

	return pn
}

func (ft funcToken) call() string {
	b := newBuffer()

	var fn string
	idx := strings.LastIndex(ft.signature, "/")

	if idx > 0 {
		fn = strings.Trim(ft.signature[idx:len(ft.signature)], "/")
	}

	// remove package prefix if same
	fn = strings.TrimPrefix(fn, ft.cpkg)
	fn = strings.TrimPrefix(fn, ".")

	args := []string{}
	for _, arg := range ft.args {
		args = append(args, arg.varName())
	}

	rets := []string{}
	for _, ret := range ft.rets {
		rets = append(rets, ret.varName())
	}

	b.ws("%s := %s(%s)", strings.Join(rets, ","), fn, strings.Join(args, ","))

	return b.buf.String()
}
