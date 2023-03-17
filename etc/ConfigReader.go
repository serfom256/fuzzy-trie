package etc

import (
	"gopkg.in/yaml.v2"
	"os"
)

func ReadConfig(configFile string) (Config, error) {
	config := &Config{}
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return *config, err
	}
	err = yaml.Unmarshal(yamlFile, config)
	return *config, err
}
