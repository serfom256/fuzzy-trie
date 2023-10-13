package core

type SearchData struct {
	Count       int
	Typos       int
	toSearch    string
	founded     []Result
	resultCache map[*TNode]bool
	nodeCache   map[*TNode]int
	onFind      OnFindFunction
}

func (data *SearchData) isFounded() bool {
	return data.Count < len(data.founded)
}

func (data *SearchData) canMoveNext(node *TNode, distance int, hash int) bool {
	return data.isFounded() || distance < 0 || node == nil || data.nodeCache[node] == hash
}
func (data *SearchData) addNodeToCache(node *TNode, hash int) {
	data.nodeCache[node] = hash
}
