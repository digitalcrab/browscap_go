// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"math"
	"sort"
	"unicode"
)

type ExpressionTree struct {
	root *node
}

func NewExpressionTree() *ExpressionTree {
	return &ExpressionTree{
		root: &node{},
	}
}

func (r *ExpressionTree) Find(userAgent []byte) string {
	res, _ := r.root.findBest(userAgent, math.MaxInt32)
	return res
}

func (r *ExpressionTree) Add(name string, lineNum int) {
	nameBytes := mapToBytes(unicode.ToLower, name)
	exp := CompileExpression(nameBytes)
	bytesPool.Put(nameBytes)

	// https://github.com/browscap/browscap/wiki/Specification:-Lookup-Algorithm:
	// If there are multiple matching patterns the longest pattern
	// should be used to identify the browser.  If there are multiple
	// patterns of the longest length then the pattern that is
	// earliest in the list of patterns should be used to identify the
	// browser.
	// The formula below implements this algorithm, with lower numbers
	// being better.
	if lineNum>>20 > 0 {
		panic("Need more bits for lineNum!")
	}
	if len(nameBytes)>>10 > 0 {
		panic("Need more bits for the UA regexp length!")
	}
	score := int32(((1<<10)-len(nameBytes))<<20 + lineNum)

	last := r.root
	for _, e := range exp {
		var found *node
		if e.Fuzzy() {
			for _, node := range last.nodesFuzzy {
				if node.token.Equal(e) {
					found = node
					break
				}
			}
		} else {
			for _, node := range last.nodesPure[e.Shard()] {
				if node.token.Equal(e) {
					found = node
					break
				}
			}
		}
		if found == nil {
			found = &node{
				token: e,
			}
			last.addChild(found)
		}
		if score < found.topScore || found.topScore == 0 {
			found.topScore = score
		}
		last = found
	}

	last.name = name
	last.score = score
}

type node struct {
	name     string
	score    int32
	topScore int32

	token Token

	nodesPure  map[byte]nodes
	nodesFuzzy nodes
}

func (n *node) addChild(a *node) {
	if a.token.Fuzzy() {
		n.nodesFuzzy = append(n.nodesFuzzy, a)
		sort.Sort(n.nodesFuzzy)
	} else {
		if n.nodesPure == nil {
			n.nodesPure = map[byte]nodes{}
		}
		shard := a.token.Shard()
		n.nodesPure[shard] = append(n.nodesPure[shard], a)
		sort.Sort(n.nodesPure[shard])
	}
}

func (n *node) findBest(s []byte, minScore int32) (res string, maxScore int32) {
	if n.topScore >= minScore {
		return "", -1
	}

	match := false
	if n.token.match != nil {
		match, s = n.token.MatchOne(s)
		if !match {
			return "", n.topScore
		}

		if n.name != "" && len(s) == 0 {
			return n.name, n.score
		}
	}

	if len(s) == 0 {
		return "", -1
	}

	for _, nd := range n.nodesPure[s[0]] {
		r, ms := nd.findBest(s, minScore)
		if ms > minScore {
			break
		}

		if r != "" {
			res = r
			minScore = ms
		}
	}

	for _, nd := range n.nodesFuzzy {
		r, ms := nd.findBest(s, minScore)
		if ms > minScore {
			break
		}

		if r != "" {
			res = r
			minScore = ms
		}
	}

	return res, minScore
}

type nodes []*node

func (n nodes) Len() int {
	return len(n)
}

func (n nodes) Less(i, j int) bool {
	return n[i].topScore < n[j].topScore
}

func (n nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
