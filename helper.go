package browscap_go

import (
	"strings"
	"regexp"
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

	rNoPrefix = regexp.MustCompile("[^a-z]")
	lenPrefix = 3
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

func getPrefix(s string) (prefix string) {
	if len(s) >= lenPrefix {
		prefix = s[0 : lenPrefix]
	} else {
		prefix = s
	}
	prefix = strings.ToLower(prefix)
	// Fallback
	if rNoPrefix.MatchString(prefix) {
		prefix = "*"
	}
	return
}
