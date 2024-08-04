package trie

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/serfom256/fuzzy-trie/trie/core"
)

func ReadDir(path string, t *core.Trie) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(path, "permission denied")
			} else {
				t.Add(info.Name(), path)
			}

			return nil
		})
	if err != nil {
		log.Println(err.Error())
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
