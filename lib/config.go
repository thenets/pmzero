package lib

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// ConfigData configs from YAML file
type ConfigData struct {
	Name string   `yaml:"name"`
	CMD  []string `yaml:"cmd"`
	CPU  struct {
		Limit int `yaml:"limit"`
	}
	Linux struct {
		User string `yaml:"user"`
	}
}

// ReadConfigFile the config file
func ReadConfigFile(filePath string) ConfigData {
	// Read config file content
	configFileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR] Config file can't be read: %v\n%v", filePath, err)
	}

	// Create config object
	t := ConfigData{}
	err = yaml.Unmarshal([]byte(configFileContent), &t)
	if err != nil {
		log.Fatalf("[ERROR] ConfigData struct can't be populated:\n%v", err)
	}
	// fmt.Printf("--- t:\n%v\n\n", t)

	return t
}

var version = "v1.0.0"
