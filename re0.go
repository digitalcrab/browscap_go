// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

const (
	TOKEN_MULTI  rune = '*'
	TOKEN_SINGLE rune = '?'

	MODE_MULTI  int = 3
	MODE_SINGLE int = 2
	MODE_STATIC int = 1
)

type Expression []*Token

type Tokens []*Token

type Token struct {
	match []byte
	skip  int
	multi bool
}

func (t *Token) String() string {
	return fmt.Sprintf("%q %d %t", string(t.match), t.skip, t.multi)
}

func (t *Token) Shard() byte {
	if len(t.match) > 0 {
		return t.match[0]
	} else {
		return 0
	}
}

func (t *Token) Equal(t1 *Token) bool {
	if !bytes.Equal(t.match, t1.match) {
		return false
	}
	if t.skip != t1.skip {
		return false
	}
	if t.multi != t1.multi {
		return false
	}
	return true
}

func (t *Token) Fuzzy() bool {
	return t.multi || t.skip > 0
}

func (t *Token) MatchOne(r []byte) (bool, []byte) {
	if t.skip > 0 {
		if len(r) < t.skip {
			return false, nil
		}
		r = r[t.skip:]
	}

	n := len(t.match)

	if n == 0 {
		return t.multi, nil
	} else if len(r) < n {
		return false, nil
	}

	if t.multi {
		ind := smallBytesIndex(r, t.match)
		if ind == -1 {
			return false, nil
		}
		return true, r[ind+n:]
	}

	if !bytes.Equal(r[:n], t.match) {
		return false, nil
	}
	return true, r[n:]
}

type parserState struct {
	lastToken *Token
	lastMode  int
	exp       Expression
}

func (ps *parserState) process(r rune) {
	mode := ps.modeByR(r)

	if ps.lastToken == nil {
		ps.lastToken = &Token{match: make([]byte, 0, 16)}
	}

	modMode := false
	if mode == MODE_MULTI || mode == MODE_SINGLE {
		modMode = true
	}

	lastModMode := false
	if ps.lastMode == MODE_MULTI || ps.lastMode == MODE_SINGLE {
		lastModMode = true
	}

	// changed
	if ps.lastMode > 0 && modMode && !lastModMode {
		ps.exp = append(ps.exp, ps.lastToken)

		ps.lastToken = &Token{match: make([]byte, 0, 16)}
	}

	// update
	switch r {
	case TOKEN_SINGLE:
		ps.lastToken.skip++
	case TOKEN_MULTI:
		ps.lastToken.multi = true
	default:
		ps.lastToken.match = appendRune(ps.lastToken.match, r)
	}

	ps.lastMode = mode
}

func (ps *parserState) modeByR(r rune) int {
	if r == TOKEN_MULTI {
		return MODE_MULTI
	} else if r == TOKEN_SINGLE {
		return MODE_SINGLE
	}
	return MODE_STATIC
}

func (ps *parserState) last() {
	// save prev token
	if ps.lastToken != nil {
		ps.exp = append(ps.exp, ps.lastToken)
	}
}

func CompileExpression(s []byte) Expression {
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

	return state.exp
}

func appendRune(b []byte, r rune) []byte {
	if r < utf8.RuneSelf {
		return append(b, byte(r))
	}

	rb := make([]byte, utf8.UTFMax)
	n := utf8.EncodeRune(rb, r)
	return append(b, rb[0:n]...)
}
