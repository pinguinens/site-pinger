package site

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseDir(dirName string) (Collection, error) {
	dn := defaultPath
	if dirName != "" {
		dn = dirName
	}

	files, err := os.ReadDir(dn)
	if err != nil {
		return Collection{}, err
	}

	sl := make([]Site, len(files), len(files))
	for i, f := range files {
		sl[i], err = ParseFile(fmt.Sprintf("%v%v%v", dn, string(os.PathSeparator), f.Name()))
		if err != nil {
			return Collection{}, err
		}
	}

	return Collection{List: sl}, nil
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
