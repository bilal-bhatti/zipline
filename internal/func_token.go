package internal

import (
	"strings"

	"github.com/bilal-bhatti/zipline/internal/tokens"
)

type funcToken struct {
	signature string
	args      []*tokens.TypeToken
	rets      []*tokens.TypeToken
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

func (ft funcToken) call(cpkgpath string) string {
	b := newBuffer()

	var fn string

	// remove package prefix if same
	fn = strings.TrimPrefix(ft.signature, cpkgpath)
	fn = strings.TrimPrefix(fn, ".")

	// if above didn't do the job
	idx := strings.LastIndex(fn, "/")

	if idx > 0 {
		fn = strings.Trim(ft.signature[idx:len(ft.signature)], "/")
	}

	args := []string{}
	for _, arg := range ft.args {
		args = append(args, arg.VarName())
	}

	rets := []string{}
	for _, ret := range ft.rets {
		rets = append(rets, ret.VarName())
	}

	b.ws("%s := %s(%s)", strings.Join(rets, ","), fn, strings.Join(args, ","))

	return b.buf.String()
}
