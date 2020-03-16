package util

import (
	"bytes"
	"fmt"
	"log"
)

// func main() {
// 	var b bytes.Buffer
// 	foo := bufio.NewWriter(&b)
// }

type Buffer struct {
	buf *bytes.Buffer
}

func NewBuffer() *Buffer {
	bites := make([]byte, 0)
	return &Buffer{
		buf: bytes.NewBuffer(bites),
	}
}

func (b *Buffer) Sprintf(s string, vals ...interface{}) {
	var err error
	if len(vals) > 0 {
		_, err = b.buf.WriteString(fmt.Sprintf(s, vals...))
	} else {
		_, err = b.buf.WriteString(s)
	}

	// if this fails no point in continuing
	if err != nil {
		log.Fatalf("writing to a byte Buffer failed with error: %v", err)
	}
}

func (b *Buffer) WriteBuffer(other *Buffer) {
	_, err := b.buf.Write(other.Bytes())
	if err != nil {
		log.Fatalf("writing to a byte Buffer failed with error: %v", err)
	}
}

func (b *Buffer) Reset() {
	b.buf.Reset()
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *Buffer) String() string {
	return b.buf.String()
}
