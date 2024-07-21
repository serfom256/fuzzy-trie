package core

type TNode struct {
	Element         byte     `json:"0"`
	Sequence        []byte   `json:"1"`
	End             bool     `json:"2"`
	SerializationId *string  `json:"3"`
	Prev            *TNode   `json:"-"`
	Successors      []*TNode `json:"4"`
	Pairs           [][]byte `json:"5"`
	SuccessorsCount int      `json:"6"`
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
	t.SuccessorsCount++
	if t.Prev != nil {
		t.Prev.SuccessorsCount++
	}
}

func (t *TNode) GetSuccessors(serializer *Serializer) []*TNode {
	t.restoreBranchIfNecessary(serializer)
	return t.Successors
}

func (t *TNode) AddPairs(pairs [][]byte, serializer *Serializer) {
	t.restoreBranchIfNecessary(serializer)
	if t.Pairs == nil {
		t.Pairs = [][]byte{}
	}
	t.Pairs = append(t.Pairs, pairs...)
}

func (t *TNode) AddPair(pair []byte, serializer *Serializer) {
	t.restoreBranchIfNecessary(serializer)
	if t.Pairs == nil {
		t.Pairs = [][]byte{}
	}
	t.Pairs = append(t.Pairs, pair)
}

func (t *TNode) IsEmpty(serializer *Serializer) bool {
	t.restoreBranchIfNecessary(serializer)
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
			t.Successors = t.removeElement(i, t.Successors, serializer)
			t.SuccessorsCount--
			if t.Prev != nil {
				t.Prev.SuccessorsCount--
			}
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

func (t *TNode) removeElement(pos int, elements []*TNode, serializer *Serializer) []*TNode {
	t.restoreBranchIfNecessary(serializer)
	length := len(elements)
	elements[pos] = elements[length-1]
	return elements[:length-1]
}

func (t *TNode) restoreBranchIfNecessary(serializer *Serializer) {
	if t.SerializationId != nil {
		t.Successors = serializer.DeserializeNode(t)
		t.restoreAfterSerialization()
		t.SerializationId = nil
	}
}
