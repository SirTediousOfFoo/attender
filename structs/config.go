// Package structs provides structures for configuration.
package structs

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Database struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Name     string `yaml:"database"`
}

type Config struct {
    Database     Database `yaml:"database"`
    PasswordSalt string   `yaml:"passwordSalt"`
}

func GetConfig(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		log.Fatalf("in file %q: %v", filename, err)
		return nil, err
	}

	return c, err
}
