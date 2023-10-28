package core

type TNode struct {
	Element         byte
	Sequence        []byte
	End             bool
	SerializationId *string
	Prev            *TNode
	Successors      []*TNode
	Pairs           [][]byte
}

func (t *TNode) FindSuccessor(b byte, serializer *Serializer) *TNode {
	t.restoreBranchIfNecessary(serializer)

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

func (t *TNode) AddSuccessor(node *TNode, serializer *Serializer) {
	t.restoreBranchIfNecessary(serializer)

	if t.Successors == nil {
		t.Successors = []*TNode{}
	}
	t.Successors = append(t.Successors, node)
}

func (t *TNode) AddPairs(pairs [][]byte) {
	if t.Pairs == nil {
		t.Pairs = [][]byte{}
	}
	t.Pairs = append(t.Pairs, pairs...)
}

func (t *TNode) AddPair(pair []byte) {
	if t.Pairs == nil {
		t.Pairs = [][]byte{}
	}
	t.Pairs = append(t.Pairs, pair)
}

func (t *TNode) IsEmpty() bool {
	return t.Sequence == nil || len(t.Sequence) == 0
}

func (t *TNode) Get(b byte, serializer *Serializer) (*TNode, bool) {
	t.restoreBranchIfNecessary(serializer)

	bt := int(b)
	for _, j := range t.Successors {
		if j.Element == b || absInt(int(j.Element)-bt) == 32 {
			return j, true
		}
	}
	return nil, false
}

func (t *TNode) RemoveSuccessor(node *TNode, serializer *Serializer) {
	t.restoreBranchIfNecessary(serializer)

	for i, element := range t.Successors {
		if element == node {
			t.Successors = removeElement(i, t.Successors)
			return
		}
	}
}

func (t *TNode) getAncestorOnDistance(distance int) *TNode {
	node := t
	for distance > 0 && node != nil {
		node = node.Prev
		distance--
	}
	return node
}

func (t *TNode) getPreviousSuccessorsSize() int {
	node := t
	size := 0

	for node != nil {
		size += len(node.Successors)
		node = node.Prev
	}
	return size
}

func (t *TNode) prepareToSerialization() {
	for _, successor := range t.Successors {
		successor.Prev = nil
		successor.prepareToSerialization()
	}
}

func (t *TNode) restoreAfterSerialization() {
	for _, successor := range t.Successors {
		successor.Prev = t
		successor.restoreAfterSerialization()
	}
}

func removeElement(pos int, elements []*TNode) []*TNode {
	length := len(elements)
	elements[pos] = elements[length-1]
	return elements[:length-1]
}

func (t *TNode) restoreBranchIfNecessary(serializer *Serializer) {
	if t.SerializationId != nil {
		serializer.deserialized++
		t.Successors = serializer.DeserializeNode(t)
		t.restoreAfterSerialization()
		t.SerializationId = nil
	}
}
