package browscap_go

import (
	"regexp"
	"bytes"
)

type expression struct {
	idx	int
	exp	*regexp.Regexp
	val	[]byte
}

func newRegexpExpression(idx int, val string) *expression {
	exp, _ := regexp.Compile("(?i)^" + escapePattern(val) + "$")
	return &expression{
		idx: idx,
		exp: exp,
	}
}

func newCompareExpression(idx int, val []byte) *expression {
	return &expression{
		idx: idx,
		val: bytes.ToLower(val),
	}
}

func (self *expression) Match(val1, val2 []byte) bool {
	if self.exp != nil {
		return self.exp.Match(val1)
	}
	if len(self.val) > 0 {
		return bytes.Equal(self.val, val2)
	}
	return false
}
