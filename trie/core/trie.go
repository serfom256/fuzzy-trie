package core

import (
	"math"
	"reflect"
	"strings"
)

type Trie struct {
	size       int
	root       *TNode
	rootNodes  map[byte]*RootNode
	serializer *Serializer
}

type RootNode struct {
	node       *TNode
	branchSize int
}

type OnFindFunction func(data SearchData, node TNode) error

const (
	ancestorDistance int = 7
)

// Add Appends a pair of key and value to map
func (t *Trie) Add(key string, value string) {
	if !checkConstraints(&key) {
		return
	}
	keyNode := t.addBranch(&key)

	if !keyNode.End {
		t.size++
	}
	keyNode.AddPair([]byte(value))
	keyNode.End = true

	//ancestor := keyNode.getAncestorOnDistance(ancestorDistance)
	//if ancestor != nil && len(key) > ancestorDistance {
	//	t.serializer.MarkNodeToBeSerialized(ancestor)
	//}
}

func (t *Trie) getNodeToSerialization(node *TNode, size int) *TNode {
	if size < ancestorDistance {
		return nil
	}
	for size > ancestorDistance {
		node = node.Prev
		size--
	}
	return node
}

func (t *Trie) Delete(key string) []string {
	node := t.findNode(key)
	if node == nil {
		return nil
	}
	node.End = false

	result := make([]string, len(node.Pairs))
	for i, val := range node.Pairs {
		result[i] = string(val[:])
	}

	return result
}

func (t *Trie) cutBranch(node *TNode) {
	if node.Successors == nil || len(node.Successors) == 0 {
		for node.Prev != nil && len(node.Prev.Successors) == 1 && !node.Prev.End {
			node = node.Prev
		}
		if node.Prev != nil {
			node.Prev.RemoveSuccessor(node, t.serializer)
			node.Prev = nil
		} else {
			delete(t.rootNodes, node.Element)
			t.root.RemoveSuccessor(node, t.serializer)
		}
	}
}
func (t *Trie) findNode(key string) *TNode {
	curr := t.root
	for i := 0; i < len(key); i++ {
		curr = curr.FindSuccessor(key[i], t.serializer)
		if curr == nil {
			return nil
		}
	}
	return curr
}

// Search returns result founded by the specified key
func (t *Trie) Search(toSearch string, distance int, cnt int, onFind OnFindFunction) []Result {
	println(t.serializer.serialized, t.serializer.deserialized)
	result := SearchData{Count: cnt, Typos: distance, toSearch: strings.ToLower(toSearch), founded: []Result{}, resultCache: map[*TNode]bool{}, nodeCache: map[*TNode]int{}, onFind: onFind}
	t.lookup(t.root, 0, distance, &result)
	return result.founded
}

func (t *Trie) lookup(curr *TNode, pos int, dist int, data *SearchData) {
	if checkHashExistence(data, pos, dist, curr) {
		return
	}
	if shouldCollectSuffix(data.toSearch, pos) {
		t.collectSuffixes(curr, data)
		return
	}
	if curr.End && isSame(t.reverseBranchLower(curr), data.toSearch, data.Typos) {
		t.collectPairs(curr, data)
	}
	if curr.Successors == nil {
		return
	}
	if pos < len(data.toSearch) {
		if next, contains := curr.Get(data.toSearch[pos], t.serializer); contains {
			t.lookup(next, pos+1, dist, data)
		}
	}
	for _, node := range curr.Successors {
		if pos < len(data.toSearch) && isCharEquals(node.Element, data.toSearch[pos]) {
			t.lookup(node, pos+1, dist, data)
		} else {
			t.lookup(node, pos+1, dist-1, data)
		}
		t.lookup(node, pos, dist-1, data)
		t.lookup(curr, pos+1, dist-1, data)
		if data.isFounded() {
			return
		}
	}
}

func (t *Trie) collectPairs(node *TNode, data *SearchData) {
	if _, ok := data.resultCache[node]; ok {
		return
	}
	if err := data.onFind(*data, *node); err != nil {
		return
	}
	data.resultCache[node] = true
	key := t.reverseBranch(node)
	var values []string
	if node.Pairs != nil {
		for _, pair := range node.Pairs {
			values = append(values, string(pair))
		}
	}
	data.founded = append(data.founded, Result{Key: key, Value: values})
}

func (t *Trie) collectSuffixes(node *TNode, data *SearchData) {
	if node == nil || data.isFounded() || data.resultCache[node] {
		return
	}
	if node.End {
		t.collectPairs(node, data)
	}
	if node.Successors == nil {
		return
	}
	for _, j := range node.Successors {
		if data.isFounded() {
			return
		}
		t.collectSuffixes(j, data)
	}
}

func (t *Trie) addBranch(key *string) *TNode {
	charArray := []byte(*key)
	rootNode := t.rootNodes[charArray[0]]

	if rootNode == nil {
		return t.insertToRoot(charArray)
	}
	curr := rootNode.node

	for i := 1; i < len(charArray); i++ {
		c := charArray[i]
		next := curr.FindSuccessor(c, t.serializer)
		if next == nil {
			seq := charArray[i:]
			if reflect.DeepEqual(seq, curr.Sequence) {
				return curr
			}
			return t.buildTree(curr, seq)
		}
		curr = next
	}
	return t.splitTree(curr)
}

func (t *Trie) splitTree(node *TNode) *TNode {
	if node.IsEmpty() {
		return node
	}
	prev := node.Prev
	if prev == nil {
		prev = t.root
	}
	prev.RemoveSuccessor(node, t.serializer)
	curr := TNode{Element: node.Element, Prev: prev}
	prev.AddSuccessor(&curr, t.serializer)
	toNext := TNode{Element: node.Sequence[0], Prev: &curr, Sequence: node.Sequence[1:]}
	toNext.End = node.End
	toNext.Pairs = node.Pairs
	node.Pairs = nil
	curr.AddSuccessor(&toNext, t.serializer)
	if prev == t.root {
		t.rootNodes[node.Element] = &RootNode{node: &curr}
	}
	return &curr
}

func (t *Trie) buildTree(node *TNode, seq []byte) *TNode {
	if node.Sequence == nil {
		newNode := TNode{Element: seq[0], Prev: node, Sequence: seq[1:]}
		node.AddSuccessor(&newNode, t.serializer)
		return &newNode
	}
	nodeSeq := node.Sequence
	node.Sequence = nil
	isEnd := node.End
	node.End = false
	tempPairs := node.Pairs
	node.Pairs = nil
	pos := 0
	length := int(math.Min(float64(len(seq)), float64(len(nodeSeq))))
	for pos < length && seq[pos] == nodeSeq[pos] {
		newNode := TNode{Element: seq[pos], Prev: node}
		node.AddSuccessor(&newNode, t.serializer)
		node = &newNode
		pos++
	}
	if pos < length {
		inserted := TNode{Element: seq[pos], Prev: node, Sequence: seq[pos+1:]}
		newNode := TNode{Element: nodeSeq[pos], Prev: node, Sequence: nodeSeq[pos+1:]}
		newNode.End = isEnd || newNode.End
		newNode.Pairs = tempPairs
		node.AddSuccessor(&newNode, t.serializer)
		node.AddSuccessor(&inserted, t.serializer)
		return &inserted
	} else if pos < len(nodeSeq) {
		newNode := TNode{Element: nodeSeq[pos], Prev: node, Sequence: nodeSeq[pos+1:]}
		newNode.End = isEnd || newNode.End
		newNode.Pairs = tempPairs
		node.AddSuccessor(&newNode, t.serializer)
		return node
	} else if pos < len(seq) {
		newNode := TNode{Element: seq[pos], Prev: node, Sequence: seq[pos+1:]}
		node.End = isEnd || node.End
		node.Pairs = tempPairs
		node.AddSuccessor(&newNode, t.serializer)
		return &newNode
	}
	return node
}

func (t *Trie) insertToRoot(key []byte) *TNode {
	first := key[0]
	rootNode := TNode{Element: first, Sequence: key[1:], End: false, Prev: nil}
	t.rootNodes[first] = &RootNode{node: &rootNode, branchSize: 1}
	t.root.AddSuccessor(&rootNode, t.serializer)
	return &rootNode
}

func checkConstraints(key *string) bool {
	if key == nil || len(*key) == 0 {
		return false
	}
	return true
}

func isCharEquals(a byte, b byte) bool {
	return absInt(int(a)-int(b)) == 32 || a == b
}

func (t *Trie) reverseBranch(node *TNode) string {
	var str []byte
	origin := node
	for node != nil {
		str = append(str, node.Element)
		node = node.Prev
	}

	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	str = append(str, origin.Sequence[:]...)
	return string(str)
}

func (t *Trie) reverseBranchLower(node *TNode) string {
	var str []byte
	origin := node
	for node != nil {
		str = append(str, node.Element)
		node = node.Prev
	}

	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	str = append(str, origin.Sequence[:]...)
	return strings.ToLower(string(str))
}

func (t *Trie) Size() int {
	return t.size
}

func InitTrie() *Trie {
	serializer := Serializer{}
	serializer.Init()
	trie := &Trie{size: 0, root: &TNode{}, rootNodes: make(map[byte]*RootNode, 32), serializer: &serializer}
	return trie
}
