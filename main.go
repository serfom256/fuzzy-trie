package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	utils "trie/etc"
	"trie/trie"
)

var config utils.Config

func main() {
	cfg, err := utils.ReadConfig("config.yaml")
	if err != nil {
		fmt.Println("An error occurred while reading config!")
	}
	config = cfg
	printConfig()
	fuzzyTrie := trie.InitTrie()
	for _, j := range config.Paths {
		readDir(j, fuzzyTrie)
	}
	start(*fuzzyTrie)
}

func search(query string, f func(string, int, int) []trie.Result) {
	t1 := time.Now()
	dist := config.Trie.Search.Distance
	fetchSize := config.Trie.Search.Fetch
	result := f(query, dist, fetchSize)
	for _, j := range result {
		fmt.Println(j.Key, j.Value)
		fmt.Println()
	}
	fmt.Println("Search time:", time.Now().Sub(t1))
}

func start(fuzzyTrie trie.Trie) {
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
	fmt.Println("Trie search config:")
	fmt.Println("\ttrie.search.distance =>", config.Trie.Search.Distance)
	fmt.Println("\ttrie.search.fetch.size =>", config.Trie.Search.Fetch)
	fmt.Println("\tpaths.to.scan =>", config.Paths)
	fmt.Println()
	fmt.Print("Indexing...")
}

func readDir(path string, t *trie.Trie) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			// readFile(path, t)
			t.Add(info.Name(), path)
			return nil
		})
	if err != nil {
		return
	}
}

func readFile(name string, t *trie.Trie) {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)
	for Scanner.Scan() {
		txt := Scanner.Text()
		t.Add(txt, file.Name())
	}
}