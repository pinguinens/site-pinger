package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigPath = "./config.yml"
)

type Config struct {
	URI    string   `yaml:"uri"`
	Domain string   `yaml:"domain"`
	Port   int      `yaml:"port"`
	Hosts  []string `yaml:"hosts"`
}

func New(path string) (*Config, error) {
	cp := path
	if path == "" {
		cp = defaultConfigPath
	}

	bytes, err := os.ReadFile(cp)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
