package trie

import trieUtils "github.com/serfom256/fuzzy-trie/etc"

type TNode struct {
	Element    byte
	Sequence   []byte
	End        bool
	Prev       *TNode
	Successors []*TNode
	Pairs      [][]byte
}

func (t *TNode) findSuccessor(b byte) *TNode {
	if t.Successors == nil {
		return nil
	}
	for _, e := range t.Successors {
		if e.Element == b {
			return e
		}
	}
	return nil
}

func (t *TNode) addSuccessor(node *TNode) {
	if t.Successors == nil {
		t.Successors = []*TNode{}
	}
	t.Successors = append(t.Successors, node)
}

func (t *TNode) addPairs(pairs [][]byte) {
	if t.Pairs == nil {
		t.Pairs = [][]byte{}
	}
	t.Pairs = append(t.Pairs, pairs...)
}

func (t *TNode) addPair(pair []byte) {
	if t.Pairs == nil {
		t.Pairs = [][]byte{}
	}
	t.Pairs = append(t.Pairs, pair)
}

func (t *TNode) isEmpty() bool {
	return t.Sequence == nil || len(t.Sequence) == 0
}

func (t *TNode) get(b byte) (*TNode, bool) {
	bt := int(b)
	for _, j := range t.Successors {
		if j.Element == b || trieUtils.AbsInt(int(j.Element)-bt) == 32 {
			return j, true
		}
	}
	return nil, false
}

func (t *TNode) removeSuccessor(node *TNode) {
	for i, element := range t.Successors {
		if element == node {
			t.Successors = removeElement(i, t.Successors)
			return
		}
	}
}

func removeElement(pos int, elements []*TNode) []*TNode {
	length := len(elements)
	elements[pos] = elements[length-1]
	return elements[:length-1]
}
