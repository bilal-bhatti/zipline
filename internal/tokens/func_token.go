package tokens

import (
	"go/token"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/util"
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

func (ft FuncToken) Call(importingPkg string, tok token.Token) string {
	buf := util.NewBuffer()

	var fn string

	// remove package prefix if same
	fn = strings.TrimPrefix(ft.Signature, importingPkg)
	fn = strings.TrimPrefix(fn, ".")

	// if above didn't do the job
	idx := strings.LastIndex(fn, "/")

	if idx > 0 {
		fn = strings.Trim(ft.Signature[idx:len(ft.Signature)], "/")
	}

	args := Join(ft.Args, func(t *TypeToken) string {
		return t.VarName()
	})

	rets := Join(ft.Rets, func(t *TypeToken) string {
		return t.VarName()
	})

	// tok is ASSN or DEFINE token (= or :=)
	buf.Sprintf("%s %s %s(%s)", rets, tok.String(), fn, args)

	return buf.String()
}
