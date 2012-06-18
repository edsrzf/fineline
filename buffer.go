package fineline

import "unicode/utf8"

// This is similar to bytes.Buffer, but it has random access
type buffer struct {
	buf       []byte
	runeBytes [utf8.UTFMax]byte // avoid allocation of slice on each WriteByte or Rune
	bootstrap [64]byte          // memory to hold first slice; helps small buffers avoid allocation.
}

func (b *buffer) grow(n int) int {
	m := len(b.buf)
	if len(b.buf)+n > cap(b.buf) {
		var buf []byte
		if b.buf == nil && n <= len(b.bootstrap) {
			buf = b.bootstrap[0:]
		} else {
			// not enough space anywhere
			buf = make([]byte, 2*cap(b.buf)+n)
			copy(buf, b.buf)
		}
		b.buf = buf
	}
	b.buf = b.buf[0 : m+n]
	return m
}

func (b *buffer) len() int {
	return len(b.buf)
}

func (b *buffer) remove(pos int) {
	if pos >= len(b.buf) {
		return
	}
	if pos < len(b.buf) {
		copy(b.buf[pos:], b.buf[pos+1:])
	}
	b.buf = b.buf[:len(b.buf)-1]
}

func (b *buffer) reset() {
	b.buf = b.buf[:0]
}

func (b *buffer) pretruncate(pos int) {
	newbuf := b.buf[pos:]
	copy(b.buf, newbuf)
	b.buf = b.buf[:len(newbuf)]
}

func (b *buffer) truncate(pos int) {
	b.buf = b.buf[:pos]
}

func (b *buffer) transpose(pos int) {
	if pos == 0 {
		return
	}
	if pos == len(b.buf) {
		pos--
	}
	b.buf[pos-1], b.buf[pos] = b.buf[pos], b.buf[pos-1]
}

func (b *buffer) Write(p []byte, pos int) {
	n := len(p)
	b.grow(n)
	copy(b.buf[pos+n:], b.buf[pos:])
	copy(b.buf[pos:], p)
}

func (b *buffer) WriteByte(c byte, pos int) error {
	b.grow(1)
	copy(b.buf[pos+1:], b.buf[pos:])
	b.buf[pos] = c
	return nil
}

func (b *buffer) WriteRune(r rune, pos int) {
	if r < utf8.RuneSelf {
		b.WriteByte(byte(r), pos)
		return
	}
	n := utf8.EncodeRune(b.runeBytes[0:], r)
	b.Write(b.runeBytes[0:n], pos)
}

func (b *buffer) WriteString(s string, pos int) {
	b.Write([]byte(s), pos)
}

func (b *buffer) Bytes() []byte {
	return b.buf
}

func (b *buffer) String() string {
	return string(b.buf)
}
