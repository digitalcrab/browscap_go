package re0

import (
	"bytes"
	"fmt"
)

const (
	TOKEN_MULTI  rune = '*'
	TOKEN_SINGLE rune = '?'

	MODE_MULTI  int = 3
	MODE_SINGLE int = 2
	MODE_STATIC int = 1
)

type Expression []iToken

type iToken interface {
}

type tokenMulti struct {
}

func (t *tokenMulti) String() string {
	return "multi: any number of any chars"
}

type tokenSingle struct {
	count int
}

func (t *tokenSingle) String() string {
	return fmt.Sprintf("single: %d of any chars", t.count)
}

type tokenStatic struct {
	buf *bytes.Buffer
}

func (t *tokenStatic) String() string {
	return fmt.Sprintf("static: %q", t.buf.String())
}

type parserState struct {
	lastToken iToken
	lastMode  int
	exp       Expression
}

func (self *parserState) process(r rune) {
	mode := self.modeByR(r)
	create := false

	// changed
	if self.lastMode != mode {

		// save prev token
		if self.lastToken != nil {
			self.exp = append(self.exp, self.lastToken)
		}

		// crate new
		create = true
	} else {
		// update
		switch r {
		case TOKEN_SINGLE:
			if self.lastToken != nil {
				self.lastToken.(*tokenSingle).count++
			} else {
				create = true
			}
		default:
			if r != TOKEN_MULTI {
				if self.lastToken != nil {
					self.lastToken.(*tokenStatic).buf.WriteRune(r)
				} else {
					create = true
				}
			}
		}
	}

	if create {
		switch r {
		case TOKEN_MULTI:
			self.lastToken = &tokenMulti{}
		case TOKEN_SINGLE:
			self.lastToken = &tokenSingle{count: 1}
		default:
			buf := new(bytes.Buffer)
			buf.WriteRune(r)
			self.lastToken = &tokenStatic{buf: buf}
		}
	}

	self.lastMode = mode
}

func (self *parserState) modeByR(r rune) int {
	if r == TOKEN_MULTI {
		return MODE_MULTI
	} else if r == TOKEN_SINGLE {
		return MODE_SINGLE
	}
	return MODE_STATIC
}

func (self *parserState) last() {
	// save prev token
	if self.lastToken != nil {
		self.exp = append(self.exp, self.lastToken)
	}
}

func (self *parserState) isMultiOrSingle(e iToken) bool {
	_, ok1 := e.(*tokenMulti)
	_, ok2 := e.(*tokenSingle)
	return ok1 || ok2
}

func (self Expression) Match(s []byte) bool {
	reader := bytes.NewBuffer(s)
	skip := false

	for _, e := range self {
		switch t := e.(type) {
		case *tokenMulti:
			skip = true
		case *tokenSingle:
			if got := reader.Next(t.count); len(got) != t.count {
				return false
			}
		case *tokenStatic:
			expect := t.buf.Bytes()
			if reader.Len() < len(expect) {
				return false
			}
			if !skip {
				got := reader.Next(len(expect))
				if !bytes.Equal(got, expect) {
					return false
				}
			} else {
				// get lat
				ind := bytes.Index(reader.Bytes(), expect)
				if ind == -1 {
					return false
				}
				reader.Next(ind + len(expect))
				skip = false
			}
		}
	}

	if !skip && reader.Len() > 0 {
		return false
	}
	return true
}

func Compile(s []byte) Expression {
	state := &parserState{
		exp:      Expression{},
		lastMode: -1,
	}
	reader := bytes.NewReader(s)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		state.process(r)
	}
	state.last()

	// check if multi and single are near to each other
checkMulti:
	for i, e := range state.exp {
		if len(state.exp)-1 > i && state.isMultiOrSingle(e) && state.isMultiOrSingle(state.exp[i+1]) {
			state.exp = append(state.exp[:i], state.exp[i+1:]...)
			goto checkMulti
		}
	}

	return state.exp
}
