package etc

import (
	"gopkg.in/yaml.v2"
	"os"
)

func ReadConfig(configFile string) Config {
	config := &Config{}
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		panic("An error occurred while reading config!")
	}
	err = yaml.Unmarshal(yamlFile, config)
	return *config
}
