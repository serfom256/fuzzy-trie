package trie

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
