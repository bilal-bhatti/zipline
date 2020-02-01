package internal

import (
	"bytes"
	"fmt"
	"strings"
)

type funcToken struct {
	cpkg, signature string
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

func (ft funcToken) call(param string) string {
	var b bytes.Buffer

	var fn string
	idx := strings.LastIndex(ft.signature, "/")

	if idx > 0 {
		fn = strings.Trim(ft.signature[idx:len(ft.signature)], "/")
	}

	// remove package prefix if same
	fn = strings.TrimPrefix(fn, ft.cpkg)
	fn = strings.TrimPrefix(fn, ".")
	b.WriteString(fn)
	b.WriteString(fmt.Sprintf("(%s)", param))

	return b.String()
}
