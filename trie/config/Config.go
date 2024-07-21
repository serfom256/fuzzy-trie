package config

type Config struct {
	Trie struct {
		Search struct {
			Distance int `yaml:"distance"`
			Fetch    int `yaml:"size"`
		} `yaml:"search"`
	} `yaml:"trie"`
	Paths []string `yaml:"paths"`
}
