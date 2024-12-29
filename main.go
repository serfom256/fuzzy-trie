package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/serfom256/fuzzy-trie/trie/config"
	"github.com/serfom256/fuzzy-trie/trie/core"
	"github.com/serfom256/fuzzy-trie/trie/core/common"
	"github.com/serfom256/fuzzy-trie/trie/core/model"
)

var context config.Config

func main() {
	context = config.ReadConfig("config.yaml")
	printConfig()
	fuzzyTrie := core.New()
	for _, j := range context.Paths {
		common.ReadDir(j, fuzzyTrie)
	}
	start(fuzzyTrie)
}

func search(query string, f func(string, int, int, core.OnFindFunction) []model.Result) {
	updateConfig()
	dist := context.Trie.Search.Distance
	fetchSize := 1
	t1 := time.Now()
	result := f(query, dist, fetchSize, func(data core.LookupResult, node core.TNode) error {
		return nil
	})

	for _, j := range result {
		fmt.Printf("\n%s = ", j.Key)
		for _, val := range j.Value {
			fmt.Printf("\n\t%s", val)
		}
		fmt.Println()
	}
	fmt.Printf("\nSearch time: %s", time.Since(t1))
}

func start(fuzzyTrie *core.Trie) {
	debug.SetGCPercent(100)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n\nIndexed total: ", fuzzyTrie.Size())
	for {
		fmt.Print("\nEnter to search: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		search(text, fuzzyTrie.Search)
	}
}

func printConfig() {
	fmt.Println("Trie search context:")
	fmt.Println("\ttrie.search.distance =>", context.Trie.Search.Distance)
	fmt.Println("\ttrie.search.fetch.size =>", context.Trie.Search.Fetch)
	fmt.Println("\tpaths.to.scan =>", context.Paths)
	fmt.Println()
	fmt.Println("Indexing...")
}

func updateConfig() {
	newConfig := config.ReadConfig("config.yaml")
	if !reflect.DeepEqual(context, newConfig) {
		context = newConfig
		fmt.Println("\nConfig updated!")
		fmt.Println()
	}
}
