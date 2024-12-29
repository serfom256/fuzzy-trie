package core

import (
	"math"
	"reflect"
	"strings"
	"sync"

	"github.com/serfom256/fuzzy-trie/trie/core/model"
)

type Trie struct {
	size       int
	root       *TNode
	rootNodes  map[byte]*RootNode
	serializer *Serializer
	scheduler  *Scheduler
	lock       sync.RWMutex
}

type RootNode struct {
	node *TNode
}

type OnFindFunction func(data LookupResult, node TNode) error

// Add Appends a pair of key and value to map
func (t *Trie) Add(key string, value string) {
	t.lock.Lock()

	if !checkConstraints(&key) {
		return
	}

	keyNode := t.addBranch(&key)

	if !keyNode.End {
		t.size++
	}
	keyNode.End = true

	keyNode.AddPair([]byte(value), t.serializer)

	t.lock.Unlock()
}

// Delete key with linked values from trie and return linked values
func (t *Trie) Delete(key string) []string {
	t.lock.Lock()

	node := t.findNode(key)
	if node == nil {
		return nil
	}
	node.End = false

	result := make([]string, len(node.Pairs))
	for i, val := range node.Pairs {
		result[i] = string(val[:])
	}
	t.cutBranch(node)

	t.lock.Unlock()
	return result
}

func (t *Trie) cutBranch(node *TNode) {
	if node.SuccessorsCount == 0 {
		for node.Prev != nil && len(node.Prev.GetSuccessors(t.serializer)) == 1 && !node.Prev.End {
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
	current := t.root
	for i := 0; i < len(key); i++ {
		prev := current
		current = current.FindSuccessor(key[i], t.serializer)
		if current == nil {
			if string(prev.Sequence) == key[i:] {
				return prev
			}
			return nil
		}
	}
	return current
}

// Search returns result founded by the specified key
func (t *Trie) Search(toSearch string, distance int, cnt int, onFind OnFindFunction) []model.Result {
	t.lock.Lock()

	result := LookupResult{Count: cnt, Typos: distance, toSearch: strings.ToLower(toSearch), founded: []model.Result{}, resultCache: map[*TNode]bool{}, nodeCache: map[*TNode]int{}, onFind: onFind}
	t.lookup(t.root, 0, distance, &result)

	t.lock.Unlock()
	return result.founded
}

func (t *Trie) lookup(current *TNode, pos int, dist int, data *LookupResult) {
	if checkHashExistence(data, pos, dist, current) {
		return
	}
	if shouldCollectSuffix(data.toSearch, pos) {
		t.collectSuffixes(current, data)
		return
	}
	if current.End && isSame(t.reverseBranchLower(current), data.toSearch, data.Typos) {
		t.collectPairs(current, data)
	}

	successors := current.GetSuccessors(t.serializer)

	if successors == nil {
		return
	}

	if pos < len(data.toSearch) {
		if next, contains := current.Get(data.toSearch[pos], t.serializer); contains {
			t.lookup(next, pos+1, dist, data)
		}
	}

	for _, node := range successors {
		if pos < len(data.toSearch) && isCharEquals(node.Element, data.toSearch[pos]) {
			t.lookup(node, pos+1, dist, data)
		} else {
			t.lookup(node, pos+1, dist-1, data)
		}
		t.lookup(node, pos, dist-1, data)
		t.lookup(current, pos+1, dist-1, data)
		if data.isFounded() {
			return
		}
	}
}

func (t *Trie) collectPairs(node *TNode, data *LookupResult) {
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
	data.founded = append(data.founded, model.Result{Key: key, Value: values})
}

func (t *Trie) collectSuffixes(node *TNode, data *LookupResult) {
	if node == nil || data.isFounded() || data.resultCache[node] {
		return
	}
	if node.End {
		t.collectPairs(node, data)
	}

	successors := node.GetSuccessors(t.serializer)

	if successors == nil {
		return
	}
	for _, j := range successors {
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
	current := rootNode.node

	for i := 1; i < len(charArray); i++ {
		c := charArray[i]
		next := current.FindSuccessor(c, t.serializer)
		if next == nil {
			seq := charArray[i:]
			if reflect.DeepEqual(seq, current.Sequence) {
				return current
			}
			return t.buildTree(current, seq)
		}
		next.Prev = current
		current = next
	}
	return t.fragmentBranch(current)
}

func (t *Trie) fragmentBranch(node *TNode) *TNode {

	if node.IsEmpty(t.serializer) {
		return node
	}

	prev := node.Prev
	if prev == nil {
		prev = t.root
	}

	prev.RemoveSuccessor(node, t.serializer)

	current := TNode{Element: node.Element, Prev: prev, SuccessorsCount: node.SuccessorsCount}
	prev.AddSuccessor(&current, t.serializer)

	nextNode := TNode{Element: node.Sequence[0], Prev: &current, Sequence: node.Sequence[1:], End: node.End, Pairs: node.Pairs}

	node.Pairs = nil // TODO relink all node successors to current
	current.AddSuccessor(&nextNode, t.serializer)

	if prev == t.root {
		t.rootNodes[current.Element] = &RootNode{node: &current}
	}
	return &current
}

func (t *Trie) buildTree(node *TNode, seq []byte) *TNode {
	if node.Sequence == nil {
		newNode := TNode{Element: seq[0], Prev: node, Sequence: seq[1:]}
		node.AddSuccessor(&newNode, t.serializer)
		return &newNode
	}

	nodeSeq := node.Sequence
	isEnd := node.End
	tempPairs := node.Pairs

	node.Sequence = nil
	node.End = false
	node.Pairs = nil
	node.SuccessorsCount = 0
	pos := 0

	length := int(math.Min(float64(len(seq)), float64(len(nodeSeq))))

	for pos < length && seq[pos] == nodeSeq[pos] {
		newNode := TNode{Element: seq[pos], Prev: node}
		node.AddSuccessor(&newNode, t.serializer)
		node = &newNode
		pos++
	}
	if pos < length {
		successor1 := TNode{Element: seq[pos], Prev: node, Sequence: seq[pos+1:]}
		successor2 := TNode{Element: nodeSeq[pos], Prev: node, Sequence: nodeSeq[pos+1:]}

		successor2.End = isEnd || successor2.End
		successor2.Pairs = tempPairs

		node.AddSuccessor(&successor2, t.serializer)
		node.AddSuccessor(&successor1, t.serializer)
		return &successor1
	} else if pos < len(nodeSeq) {
		newNode := TNode{Element: nodeSeq[pos], Prev: node, Sequence: nodeSeq[pos+1:], SuccessorsCount: node.SuccessorsCount}
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
	t.rootNodes[first] = &RootNode{node: &rootNode}
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

func New() *Trie {
	serializer := Serializer{}
	serializer.Init()

	trie := &Trie{
		size:       0,
		root:       &TNode{},
		rootNodes:  make(map[byte]*RootNode, 32),
		serializer: &serializer,
	}
	trie.scheduler = NewScheduler(trie)

	return trie
}
