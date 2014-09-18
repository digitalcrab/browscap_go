package browscap_go

import (
	"strings"
	"bytes"
)

var (
	regexReplace = map[string]string{
		// Basic
		"(": "\\(",
		")": "\\)",
		"[": "\\[",
		"]": "\\]",
		"{": "\\{",
		"}": "\\}",
		"<": "\\<",
		">": "\\>",
		"$": "\\$",
		"^": "\\^",
		"+": "\\+",
		"!": "\\!",
		"=": "\\=",
		"|": "\\|",
		":": "\\:",
		"-": "\\-",
		// Search
		"*": ".*",
		"?": ".",
	}
)

func escapePattern(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, ".", "\\.", -1)

	for k, v := range regexReplace {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

func inList(val []byte, list[][]byte) bool {
	for _, v := range list {
		if bytes.Equal(val, v) {
			return true
		}
	}
	return false
}
