package trie

type SearchData struct {
	count       int
	typos       int
	toSearch    string
	founded     []Result
	resultCache map[*TNode]bool
	nodeCache   map[*TNode]int
}

func (data *SearchData) isFounded() bool {
	return data.count < len(data.founded)
}
