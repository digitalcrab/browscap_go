package browscap_go

import "strings"

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

