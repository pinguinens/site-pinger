package site

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

const (
	defaultPath = "./sites"
)

type Site struct {
	Target Target `yaml:"target"`
}

func ParseDir(dirName string) ([]Site, error) {
	dn := defaultPath
	if dirName != "" {
		dn = dirName
	}

	files, err := os.ReadDir(dn)
	if err != nil {
		return nil, err
	}

	sl := make([]Site, len(files), len(files))
	for i, f := range files {
		sl[i], err = ParseFile(fmt.Sprintf("%v%v%v", dn, string(os.PathSeparator), f.Name()))
		if err != nil {
			return nil, err
		}
	}

	return sl, nil
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
