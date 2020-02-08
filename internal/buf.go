package internal

import (
	"bytes"
	"fmt"
)

type buffer struct {
	buf *bytes.Buffer
}

func newBuffer() buffer {
	bites := make([]byte, 0)
	return buffer{
		buf: bytes.NewBuffer(bites),
	}
}

func (b buffer) ws(s string, vals ...interface{}) {
	if len(vals) > 0 {
		b.buf.WriteString(fmt.Sprintf(s, vals...))
	} else {
		b.buf.WriteString(s)
	}
}
