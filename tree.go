// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
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
	if len(userAgent) == 0 {
		return ""
	}

	res, _ := r.root.findBest(userAgent, 0, 0)
	return res
}

func (r *ExpressionTree) Add(name string, lineNum int) {
	nameBytes := mapToBytes(unicode.ToLower, name)
	exp := CompileExpression(nameBytes)
	bytesPool.Put(nameBytes)

	score := len(name)

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
				token:    e,
				topScore: score,
			}
			last.addChild(found)
		} else if score > found.topScore {
			found.topScore = score
		}
		last = found
	}

	last.name = name
}

type node struct {
	name string

	token Token

	nodesPure  map[byte]nodes
	nodesFuzzy nodes

	topScore int
}

func (n *node) addChild(a *node) {
	if a.token.Fuzzy() {
		n.nodesFuzzy = append(n.nodesFuzzy, a)
		sort.Sort(sort.Reverse(n.nodesFuzzy))
	} else {
		if n.nodesPure == nil {
			n.nodesPure = map[byte]nodes{}
		}
		shard := a.token.Shard()
		n.nodesPure[shard] = append(n.nodesPure[shard], a)
		sort.Sort(sort.Reverse(n.nodesPure[shard]))
	}
}

func (n *node) findBest(s []byte, minScore int, x int) (res string, maxScore int) {
	if n.token.match != nil {
		match := false
		match, s = n.token.MatchOne(s)
		if !match {
			return "", n.topScore
		}

		if len(s) == 0 && len(n.nodesFuzzy) == 0 {
			return n.name, n.topScore
		}
	}

	if len(s) > 0 {
		for _, nd := range n.nodesPure[s[0]] {
			if nd.topScore <= minScore {
				break
			}
			r, ms := nd.findBest(s, minScore, x+1)
			if r != "" && ms > minScore {
				res = r
				minScore = ms
			}
		}
	}

	for _, nd := range n.nodesFuzzy {
		if nd.topScore <= minScore {
			break
		}
		r, ms := nd.findBest(s, minScore, x+1)
		if r != "" && ms > minScore {
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
