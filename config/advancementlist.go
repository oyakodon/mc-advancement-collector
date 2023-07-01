package config

import (
	"os"

	"com.oykdn.mc-advancement-collector/model"
	"gopkg.in/yaml.v2"
)

type AdvancementRecord struct {
	Criteria    []string              `yaml:"criteria"`
	Parent      string                `yaml:"parent"`
	LanguageKey string                `yaml:"languageKey"`
	Calculate   model.CalculateType   `yaml:"calculate"`
	Hidden      bool                  `yaml:"hidden"`
	Type        model.AdvancementType `yaml:"type"`
	Icon        AdvancementRecordIcon `yaml:"icon"`
}

type AdvancementRecordIcon struct {
	Url       string `yaml:"url"`
	InvSprite bool   `yaml:"invsprite"`
	Pos       int    `yaml:"pos"`
}

type AdvancementList struct {
	Advancements map[string]AdvancementRecord `yaml:"advancements"`
}

func LoadAdvancementList(path string) (*AdvancementList, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var list AdvancementList
	if err := yaml.Unmarshal(b, &list); err != nil {
		return nil, err
	}

	return &list, err
}
