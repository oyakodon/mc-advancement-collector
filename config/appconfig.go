package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	AdvancementPath string         `yaml:"advancementPath"`
	Language        string         `yaml:"language"`
	Cache           int            `yaml:"cache"`
	Assets          AppConfigAsset `yaml:"assets"`
}

type AppConfigAsset struct {
	Background map[string]AppConfigAssetBackground `yaml:"background"`
}

type AppConfigAssetBackground struct {
	Incomplete string `yaml:"incomplete"`
	Completed  string `yaml:"completed"`
}

func LoadAppConfig(path string) (*AppConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf AppConfig
	if err := yaml.Unmarshal(b, &conf); err != nil {
		return nil, err
	}

	return &conf, err
}
