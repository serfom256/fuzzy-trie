package trie

import (
	"math"
	"reflect"
	"strings"
	"unicode"

	"github.com/agnivade/levenshtein"
)

type Trie struct {
	size      int
	root      *TNode
	rootNodes map[byte]*TNode
}

const (
	suffixRegex = '*'
)

// Add Appends a pair of key and value to map
func (t *Trie) Add(key string, value string) {
	checkConstraints(&key)
	keyNode := t.addSequence(&key)
	if !keyNode.end {
		t.size++
	}
	keyNode.addPair([]byte(value))
	keyNode.end = true
}

// Search returns result founded by the specified key
func (t *Trie) Search(toSearch string, distance int, cnt int) []Result {
	result := SearchData{count: cnt, typos: distance, toSearch: strings.ToLower(toSearch), founded: []Result{}, cache: map[*TNode]bool{}}
	t.lookup(t.root, 0, distance, &result)
	return result.founded
}

func (t *Trie) lookup(curr *TNode, pos int, typos int, data *SearchData) {
	if typos < 0 || curr == nil || data.isFounded() {
		return
	}
	if pos < len(data.toSearch) && data.toSearch[pos] == suffixRegex {
		t.collectSuffixes(curr, data)
		return
	}
	if curr.end && isSame(t.reverseBranchLower(curr), data.toSearch, data.typos) { //todo optimize with with caching
		t.collectPairs(curr, data)
	}
	if curr.successors == nil {
		return
	}
	hasNext := false
	if pos < len(data.toSearch) {
		if next, contains := curr.get(data.toSearch[pos]); contains {
			t.lookup(next, pos+1, typos, data)
			hasNext = true
		}
	}
	for _, node := range curr.successors {
		if !hasNext && pos < len(data.toSearch) && isCharEquals(node.element, data.toSearch[pos]) {
			t.lookup(node, pos+1, typos, data)
		} else {
			t.lookup(node, pos+1, typos-1, data)
		}
		t.lookup(node, pos, typos-1, data)
		t.lookup(curr, pos+1, typos-1, data)
		if data.isFounded() {
			return
		}
	}
}

func (t *Trie) collectPairs(node *TNode, data *SearchData) {
	if _, ok := data.cache[node]; ok {
		return
	}
	data.cache[node] = true
	key := t.reverseBranch(node)
	var values []string
	if node.pairs != nil {
		for _, pair := range node.pairs {
			values = append(values, string(pair))
		}
	}
	data.founded = append(data.founded, Result{Key: key, Value: values})
}

func (t *Trie) collectSuffixes(node *TNode, data *SearchData) {
	if node == nil || data.isFounded() || data.cache[node] {
		return
	}
	if node.end {
		t.collectPairs(node, data)
	}
	if node.successors == nil {
		return
	}
	for _, j := range node.successors {
		if data.isFounded() {
			return
		}
		t.collectSuffixes(j, data)
	}
}

func (t *Trie) addSequence(key *string) *TNode {
	charArray := []byte(*key)
	curr := t.rootNodes[charArray[0]]
	if curr == nil {
		return t.insertToRoot(charArray)
	}
	for i := 1; i < len(charArray); i++ {
		c := charArray[i]
		next := curr.findSuccessor(c)
		if next == nil {
			seq := charArray[i:]
			if reflect.DeepEqual(seq, curr.sequence) {
				return curr
			}
			return t.buildTree(curr, seq)
		}
		curr = next
	}
	return t.splitTree(curr)
}

func isCharEquals(a byte, b byte) bool {
	return math.Abs(float64(int(a)-int(b))) == 32 || a == b
}

func (t *Trie) splitTree(node *TNode) *TNode {
	if node.isEmpty() {
		return node
	}
	prev := node.prev
	if prev == nil {
		prev = t.root
	}
	prev.removeSuccessor(node)
	curr := TNode{element: node.element, prev: prev}
	prev.addSuccessor(&curr)
	toNext := TNode{element: node.sequence[0], prev: &curr, sequence: node.sequence[1:]}
	toNext.end = node.end
	toNext.pairs = node.pairs
	node.pairs = nil
	curr.addSuccessor(&toNext)
	if prev == t.root {
		t.rootNodes[node.element] = &curr
	}
	return &curr
}

func (t *Trie) insertToRoot(key []byte) *TNode {
	first := key[0]
	rootNode := TNode{element: first, sequence: key[1:], end: false, prev: nil}
	t.rootNodes[first] = &rootNode
	t.root.addSuccessor(&rootNode)
	return &rootNode
}

func checkConstraints(key *string) {
	if key == nil || len(*key) == 0 {
		panic("Specified Key is empty")
	}
}

func (t *Trie) buildTree(node *TNode, seq []byte) *TNode {
	if node.sequence == nil {
		newNode := TNode{element: seq[0], prev: node, sequence: seq[1:]}
		node.addSuccessor(&newNode)
		return &newNode
	}
	nodeSeq := node.sequence
	node.sequence = nil
	isEnd := node.end
	node.end = false
	tempPairs := node.pairs
	node.pairs = nil
	pos := 0
	length := int(math.Min(float64(len(seq)), float64(len(nodeSeq))))
	for pos < length && seq[pos] == nodeSeq[pos] {
		newNode := TNode{element: seq[pos], prev: node}
		node.addSuccessor(&newNode)
		node = &newNode
		pos++
	}
	if pos < length {
		inserted := TNode{element: seq[pos], prev: node, sequence: seq[pos+1:]}
		newNode := TNode{element: nodeSeq[pos], prev: node, sequence: nodeSeq[pos+1:]}
		newNode.end = isEnd || newNode.end
		newNode.pairs = tempPairs
		node.addSuccessor(&newNode)
		node.addSuccessor(&inserted)
		return &inserted
	} else if pos < len(nodeSeq) {
		newNode := TNode{element: nodeSeq[pos], prev: node, sequence: nodeSeq[pos+1:]}
		newNode.end = isEnd || newNode.end
		newNode.pairs = tempPairs
		node.addSuccessor(&newNode)
		return node
	} else if pos < len(seq) {
		newNode := TNode{element: seq[pos], prev: node, sequence: seq[pos+1:]}
		node.end = isEnd || node.end
		node.pairs = tempPairs
		node.addSuccessor(&newNode)
		return &newNode
	}
	return node
}

func (t *Trie) reverseBranch(node *TNode) string {
	var str []byte
	origin := node
	for node != nil {
		str = append(str, node.element)
		node = node.prev
	}

	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	str = append(str, origin.sequence[:]...)
	return string(str)
}

func (t *Trie) reverseBranchLower(node *TNode) string {
	var str []byte
	origin := node
	for node != nil {
		str = append(str, byte(unicode.ToLower(rune(node.element))))
		node = node.prev
	}

	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	str = append(str, origin.sequence[:]...)
	return string(str)
}

func isSame(s1 string, s2 string, distance int) bool {
	return levenshtein.ComputeDistance(s1, s2) <= distance
}

func (t *Trie) Print() {

}

func (t *Trie) Size() int {
	return t.size
}

func InitTrie() *Trie {
	trie := Trie{size: 0, root: &TNode{}, rootNodes: make(map[byte]*TNode, 32)}
	return &trie
}
