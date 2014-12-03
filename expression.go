package browscap_go

import (
	"bytes"
	"github.com/fromYukki/browscap_go/re0"
)

type expression struct {
	Name string
	exp  re0.Expression
	val  []byte
}

func newRegexpExpression(val string) *expression {
	return &expression{
		Name: val,
		exp:  re0.Compile(bytes.ToLower([]byte(val))),
	}
}

func newCompareExpression(val string) *expression {
	return &expression{
		Name: val,
		val:  bytes.ToLower([]byte(val)),
	}
}

func (self *expression) Match(val []byte) bool {
	if self.exp != nil {
		return self.exp.Match(val)
	}
	if len(self.val) > 0 {
		return bytes.Equal(self.val, val)
	}
	return false
}

type expressionByNameLen []*expression

func (el expressionByNameLen) Len() int {
	return len(el)
}

func (el expressionByNameLen) Less(i, j int) bool {
	return len(el[i].Name) > len(el[j].Name)
}

func (el expressionByNameLen) Swap(i, j int) {
	el[i], el[j] = el[j], el[i]
}
