package trie

type SearchData struct {
	count    int
	typos    int
	toSearch string
	founded  []Result
	cache    map[*TNode]bool
}

func (data *SearchData) isFounded() bool {
	return data.count < len(data.founded)
}
