package browscap_go

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	rNoPrefix = regexp.MustCompile("[^a-z]")
	lenPrefix = 3
)

func inList(val []byte, list [][]byte) bool {
	for _, v := range list {
		if bytes.Equal(val, v) {
			return true
		}
	}
	return false
}

func getPrefix(s string) (prefix string) {
	if len(s) >= lenPrefix {
		prefix = s[0:lenPrefix]
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
