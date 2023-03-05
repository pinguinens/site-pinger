package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

const (
	defaultConfigPath = "./config.yml"
)

type Config struct {
	LogFileName   string `yaml:"log_file"`
	ConsoleFormat string `yaml:"console_format" default:"pretty"`
	SiteDir       string `yaml:"site_dir" default:"./sites"`
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

	err = c.setDefaults()
	if err != nil {
		return nil, err
	}

	err = c.validate()
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) setDefaults() error {
	return defaults.Set(c)
}

func (c *Config) validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}

	return nil
}
