package pb

import (
	"crypto/md5"
	"io"
)

// Hash ...
func (x *Req) Hash() [16]byte {

	m := md5.New()

	io.WriteString(m, x.Method)
	m.Write([]byte{0x00})
	io.WriteString(m, x.Method)
	m.Write([]byte{0x00})
	m.Write(x.Body)

	var h [16]byte
	copy(h[:], m.Sum(nil)[:])
	return h
}
