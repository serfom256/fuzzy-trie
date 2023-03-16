package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	trie "trie/trie"
)

func main() {
	x := trie.InitTrie()
	readDir("/", x)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Indexed total: ", x.Size())
	for {
		fmt.Print("Enter to search: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		search(text[:len(text)-1], x.Search)

	}
}
func search(query string, f func(string, int, int) map[string][]string) {
	t1 := time.Now()
	result := f(query, 2, 100)
	for j, i := range result {
		fmt.Println(j, i)
		fmt.Println("\n\n")
	}
	fmt.Println(time.Now().Sub(t1))
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
