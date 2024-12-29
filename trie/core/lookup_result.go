package core

import "github.com/serfom256/fuzzy-trie/trie/core/model"

type LookupResult struct {
	Count       int
	Typos       int
	toSearch    string
	founded     []model.Result
	resultCache map[*TNode]bool
	nodeCache   map[*TNode]int
	onFind      OnFindFunction
}

func (data *LookupResult) isFounded() bool {
	return data.Count < len(data.founded)
}

func (data *LookupResult) canMoveNext(node *TNode, distance int, hash int) bool {
	return data.isFounded() || distance < 0 || node == nil || data.nodeCache[node] == hash
}
func (data *LookupResult) addNodeToCache(node *TNode, hash int) {
	data.nodeCache[node] = hash
}
