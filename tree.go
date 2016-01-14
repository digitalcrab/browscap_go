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
	res, _ := r.root.findBest(userAgent, 0)
	return res
}

func (r *ExpressionTree) Add(name string) {
	nameBytes := mapToBytes(unicode.ToLower, name)
	defer bytesPool.Put(nameBytes)

	exp := CompileExpression(nameBytes)

	last := r.root
	for _, e := range exp {
		shard := e.Shard()

		var found *node
		for _, node := range last.nodesPure[shard] {
			if node.token.Equal(e) {
				found = node
				break
			}
		}
		if found == nil {
			for _, node := range last.nodesFuzzy {
				if node.token.Equal(e) {
					found = node
					break
				}
			}
		}
		if found == nil {
			found = &node{
				token:  e,
				parent: last,
			}
			if e.Fuzzy() {
				last.nodesFuzzy = append(last.nodesFuzzy, found)
				sort.Sort(last.nodesFuzzy)
			} else {
				if last.nodesPure == nil {
					last.nodesPure = map[byte]nodes{}
				}
				last.nodesPure[shard] = append(last.nodesPure[shard], found)
				sort.Sort(last.nodesPure[shard])
			}
		}
		last = found
	}

	score := len(name)

	last.name = name
	last.score = score

	for last != nil {
		if score > last.topScore {
			last.topScore = score
		}
		last = last.parent
	}

}

type node struct {
	name  string
	score int

	nodesPure  map[byte]nodes
	nodesFuzzy nodes

	token    *Token
	topScore int
	parent   *node

	lastMatch *node
}

func (n *node) findBest(s []byte, minScore int) (res string, maxScore int) {
	if n.topScore <= minScore {
		return "", -1
	}

	match := false
	if n.token != nil {
		match, s = n.token.MatchOne(s)
		if !match {
			return "", n.topScore
		}

		if n.name != "" {
			res = n.name
			minScore = n.score
		}
	}

	if len(s) > 0 {
		if n.lastMatch != nil {
			r, ms := n.lastMatch.findBest(s, minScore)
			if r != "" && ms > minScore {
				res = r
				minScore = ms
			}
		}

		for i, nd := range n.nodesPure[s[0]] {
			if nd == n.lastMatch {
				continue
			}

			r, ms := nd.findBest(s, minScore)
			if ms < minScore {
				break
			}

			if r != "" {
				res = r
				minScore = ms
				if i > 0 {
					n.lastMatch = nd
				}
			}
		}

		for i, nd := range n.nodesFuzzy {
			if nd == n.lastMatch {
				continue
			}

			r, ms := nd.findBest(s, minScore)
			if ms < minScore {
				break
			}

			if r != "" {
				res = r
				minScore = ms
				if i > 0 {
					n.lastMatch = nd
				}
			}
		}
	}

	return res, minScore
}

type nodes []*node

func (n nodes) Len() int {
	return len(n)
}

func (n nodes) Less(i, j int) bool {
	// Sort reverse
	return n[i].topScore > n[j].topScore
}

func (n nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
