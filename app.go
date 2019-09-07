package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// ConfigData configs from YAML file
type ConfigData struct {
    Name string `yaml:"name"`
    CMD []string `yaml:"cmd"`
	CPU  struct {
		Limit int `yaml:"limit"`
	}
	Linux struct {
		User string `yaml:"user"`
	}
}

func main() {
	myConfigData, err := ioutil.ReadFile("./samples/sleep.yaml")
	if err != nil {
        log.Fatalf("error: %v", err)
    }
	fmt.Print(string(myConfigData))

	t := ConfigData{}

	err = yaml.Unmarshal([]byte(myConfigData), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)
}
