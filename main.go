package main

import (
	"bufio"
	"fmt"
	"github.com/serfom256/fuzzy-trie/trie"
	"github.com/serfom256/fuzzy-trie/trie/config"
	"github.com/serfom256/fuzzy-trie/trie/core"
	"os"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

var context config.Config

func main() {
	context = config.ReadConfig("config.yaml")
	printConfig()
	fuzzyTrie := core.InitTrie()
	for _, j := range context.Paths {
		trie.ReadDir(j, fuzzyTrie)
	}
	start(fuzzyTrie)
}

func search(query string, f func(string, int, int, core.OnFindFunction) []core.Result) {
	updateConfig()
	dist := context.Trie.Search.Distance
	fetchSize := context.Trie.Search.Fetch
	t1 := time.Now()
	result := f(query, dist, fetchSize, func(data core.SearchData, node core.TNode) error {
		return nil
	})
	for _, j := range result {
		fmt.Println(j.Key, j.Value)
		fmt.Println()
	}
	fmt.Println("Search time:", time.Now().Sub(t1))
}

func start(fuzzyTrie *core.Trie) {
	debug.SetGCPercent(100)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n\nIndexed total: ", fuzzyTrie.Size())
	for {
		fmt.Print("\nEnter to search: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		search(text[:len(text)-1], fuzzyTrie.Search)
	}
}

func printConfig() {
	fmt.Println("Trie search context:")
	fmt.Println("\ttrie.search.distance =>", context.Trie.Search.Distance)
	fmt.Println("\ttrie.search.fetch.size =>", context.Trie.Search.Fetch)
	fmt.Println("\tpaths.to.scan =>", context.Paths)
	fmt.Println()
	fmt.Print("Indexing...")
}

func updateConfig() {
	newConfig := config.ReadConfig("config.yaml")
	if !reflect.DeepEqual(context, newConfig) {
		context = newConfig
		fmt.Println("\nConfig updated!")
		fmt.Println()
	}
}
