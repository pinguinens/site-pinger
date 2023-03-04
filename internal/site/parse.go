package site

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	fileSuffixS = ".yml"
	fileSuffixF = ".yaml"
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

	var count uint16
	for _, f := range files {
		if strings.HasSuffix(f.Name(), fileSuffixS) || strings.HasSuffix(f.Name(), fileSuffixF) {
			count++
		}
	}

	sl := make([]Site, 0, count)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), fileSuffixS) || strings.HasSuffix(f.Name(), fileSuffixF) {
			s, err := ParseFile(fmt.Sprintf("%v%v%v", dn, string(os.PathSeparator), f.Name()))
			if err != nil {
				return Collection{}, err
			}
			sl = append(sl, s)
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
