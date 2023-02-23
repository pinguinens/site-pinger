package site

import (
	"os"

	"gopkg.in/yaml.v3"
	
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

type Site struct {
	Target Target `yaml:"target"`
}

func ParseDir(dirName string) ([]Site, error) {
	return nil, nil
}

func ParseFile(fileName string) (Site, error) {
	s := Site{}

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return s, err
	}

	err = yaml.Unmarshal(bytes, &s)
	if err != nil {
		return s, err
	}

	err = s.setDefaults()
	if err != nil {
		return s, err
	}

	err = s.validate()
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s *Site) setDefaults() error {
	return defaults.Set(s)
}

func (s *Site) validate() error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		return err
	}

	return nil
}
