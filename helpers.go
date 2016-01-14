// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"sync"
	"unicode/utf8"
)

var (
	bytesPool = &sync.Pool{}
	minCap    = 128
)

func getBytes(size int) []byte {
	if b := bytesPool.Get(); b != nil {
		bs := b.([]byte)
		if cap(bs) >= size {
			return bs[:size]
		}
	}
	c := size
	if c < minCap {
		c = minCap
	}
	return make([]byte, size, c)
}

func mapToBytes(mapping func(rune) rune, s string) []byte {
	// In the worst case, the string can grow when mapped, making
	// things unpleasant.  But it's so rare we barge in assuming it's
	// fine.  It could also shrink but that falls out naturally.
	maxbytes := len(s) // length of b
	nbytes := 0        // number of bytes encoded in b
	// The output buffer b is initialized on demand, the first
	// time a character differs.
	var b []byte

	for i, c := range s {
		r := mapping(c)
		if b == nil {
			if r == c {
				continue
			}
			b = getBytes(maxbytes)
			nbytes = copy(b, s[:i])
		}
		if r >= 0 {
			wid := 1
			if r >= utf8.RuneSelf {
				wid = utf8.RuneLen(r)
			}
			if nbytes+wid > maxbytes {
				// Grow the buffer.
				maxbytes = maxbytes*2 + utf8.UTFMax
				nb := getBytes(maxbytes)
				copy(nb, b[0:nbytes])
				b = nb
			}
			nbytes += utf8.EncodeRune(b[nbytes:maxbytes], r)
		}
	}
	if b == nil {
		b = getBytes(maxbytes)
		copy(b, s)
		return b
	}
	return b[0:nbytes]
}

func smallBytesIndex(s, sep []byte) int {
	c := sep[0]

start:
	for i, ch := range s[:len(s)-len(sep)+1] {
		if c != ch {
			continue
		}

		for j := len(sep) - 1; j > 0; j-- {
			if s[i+j] != sep[j] {
				continue start
			}
		}
		return i
	}
	return -1
}
