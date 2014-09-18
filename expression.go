package browscap_go

import (
	"regexp"
	"bytes"
)

type expression struct {
	Name	string
	exp		*regexp.Regexp
	val		[]byte
}

func newRegexpExpression(val string) *expression {
	exp, _ := regexp.Compile("(?i)^" + escapePattern(val) + "$")
	return &expression{
		Name:	val,
		exp:	exp,
	}
}

func newCompareExpression(val string) *expression {
	return &expression{
		Name:	val,
		val:	bytes.ToLower([]byte(val)),
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
