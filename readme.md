## Fuzzy search trie

## About

This is the radix trie that optimized for fuzzy search and fuzzy suggestion by prefix

---

## Build

`go build -o bin/trie main.go`

---

## API Usage

- Add (key, value)
- Search (key, fuzzy search distance, the number of results to be fetched)

```go
fuzzyTrie := trie.InitTrie()

fuzzyTrie.Add("key", "value")
fuzzyTrie.Add("key1", "value2")
fuzzyTrie.Add("__key__", "value3")
fuzzyTrie.Add("-key-", "value4")
fuzzyTrie.Add("key__5", "value4")

fuzzyTrie.Search("key", 1, 10)
// returns [{key [value]} {key1 [value2]}]

fuzzyTrie.Search("key", 2, 10)
// returns [{key [value]} {key1 [value2]} {-key- [value4]}]

fuzzyTrie.Search("key*", 0, 10)
// returns [{key [value]} {key1 [value2]} {key__5 [value4]}]

fuzzyTrie.Search("*", 0, 10)
// returns [{key [value]} {key1 [value2]} {key__5 [value4]} {__key__ [value3]} {-key- [value4]}

fuzzyTrie.Search("=key_", 3, 10)
// returns [{__key__ [value3]} {key [value]} {key__5 [value4]} {key1 [value2]} {-key- [value4]}]

```

#### Symbol `*` means match all suffixes

---

### Features

- Support case-insensitive fuzzy search and fuzzy matching by prefix

---

## License

#### MIT, check the LICENSE file.