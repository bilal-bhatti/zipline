package internal

import (
	"bytes"
	"fmt"
	"log"
)

type buffer struct {
	buf *bytes.Buffer
}

func newBuffer() *buffer {
	bites := make([]byte, 0)
	return &buffer{
		buf: bytes.NewBuffer(bites),
	}
}

func (b *buffer) ws(s string, vals ...interface{}) {
	var err error
	if len(vals) > 0 {
		_, err = b.buf.WriteString(fmt.Sprintf(s, vals...))
	} else {
		_, err = b.buf.WriteString(s)
	}

	// if this fails no point in continuing
	if err != nil {
		log.Fatalf("writing to a byte buffer failed with error: %v", err)
	}
}

func (b *buffer) add(other *buffer) {
	_, err := b.buf.Write(other.buf.Bytes())
	if err != nil {
		log.Fatalf("writing to a byte buffer failed with error: %v", err)
	}
}
