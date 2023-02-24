package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

const (
	defaultConfigPath = "./config.yml"
	timeLayout        = "2006-01-02T15-04-05"
)

type Config struct {
	LogFile string `yaml:"log_file" default:"./d_%v.log"`
	SiteDir string `yaml:"site_dir" default:"./sites"`

	DialerTimeout   time.Duration `yaml:"dialer_timeout" default:"30"`
	DialerKeepAlive time.Duration `yaml:"dialer_keepalive" default:"30"`
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

func (c *Config) GetLogFile() string {
	lf := fmt.Sprintf(c.LogFile, time.Now().Format(timeLayout))

	return lf
}
