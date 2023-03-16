package trie

type TNode struct {
	element    byte
	sequence   []byte
	end        bool
	prev       *TNode
	successors []*TNode
	pairs      [][]byte
}

func (t *TNode) findSuccessor(b byte) *TNode {
	if t.successors == nil {
		return nil
	}
	for _, e := range t.successors {
		if e.element == b {
			return e
		}
	}
	return nil
}
func (t *TNode) addSuccessor(node *TNode) {
	if t.successors == nil {
		t.successors = []*TNode{}
	}
	t.successors = append(t.successors, node)
}

func (t *TNode) addPairs(pairs [][]byte) {
	if t.pairs == nil {
		t.pairs = [][]byte{}
	}
	t.pairs = append(t.pairs, pairs...)
}
func (t *TNode) addPair(pair []byte) {
	if t.pairs == nil {
		t.pairs = [][]byte{}
	}
	t.pairs = append(t.pairs, pair)
}

func (t *TNode) isEmpty() bool {
	return t.sequence == nil || len(t.sequence) == 0
}

func (t *TNode) isRoot() bool {
	return t.prev == nil
}

func (t *TNode) removeSuccessor(node *TNode) {
	for i, element := range t.successors {
		if element == node {
			t.successors = removeElement(i, t.successors)
			return
		}
	}
}

func removeElement(pos int, elements []*TNode) []*TNode {
	length := len(elements)
	elements[pos] = elements[length-1]
	return elements[:length-1]
}
