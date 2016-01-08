package browscap_go

import "unicode/utf8"

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
			b = make([]byte, maxbytes)
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
				nb := make([]byte, maxbytes)
				copy(nb, b[0:nbytes])
				b = nb
			}
			nbytes += utf8.EncodeRune(b[nbytes:maxbytes], r)
		}
	}
	if b == nil {
		return []byte(s)
	}
	return b[0:nbytes]
}
