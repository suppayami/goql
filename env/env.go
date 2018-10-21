package env

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Env Environment info
type Env struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// ReadEnv read environment info from env.yaml
func (e *Env) ReadEnv() *Env {

	yamlFile, err := ioutil.ReadFile("env.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, e)

	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return e
}
