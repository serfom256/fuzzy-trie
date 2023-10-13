package trie

import (
	"bufio"
	"github.com/serfom256/fuzzy-trie/trie/core"
	"os"
	"path/filepath"
)

func ReadDir(path string, t *core.Trie) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			t.Add(info.Name(), path)
			t.Add(path, path)
			return nil
		})
	if err != nil {
		return
	}
}

func ReadFile(name string, t *core.Trie) {
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
