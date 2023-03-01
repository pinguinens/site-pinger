package site

import (
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

const (
	defaultPath = "./sites"
)

type Site struct {
	Target Target `yaml:"target"`
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
