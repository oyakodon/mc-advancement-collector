package lang

import (
	"encoding/json"
	"os"
)

const (
	LANG_SUFFIX_TITLE       = ".title"
	LANG_SUFFIX_DESCRIPTION = ".description"
)

type Lang struct {
	Mapping map[string]string
}

func LoadLang(path string) (*Lang, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var mapping map[string]string
	if err := json.Unmarshal(b, &mapping); err != nil {
		return nil, err
	}

	return &Lang{
		Mapping: mapping,
	}, err
}
