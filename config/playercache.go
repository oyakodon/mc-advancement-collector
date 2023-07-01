package config

import (
	"os"

	"com.oykdn.mc-advancement-collector/model"
	"gopkg.in/yaml.v2"
)

type PlayerCache struct {
	Players map[string]model.PlayerProfile `json:"players"`
}

func LoadPlayerCache(path string) (*PlayerCache, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return &PlayerCache{
			Players: make(map[string]model.PlayerProfile),
		}, nil
	}

	var p PlayerCache
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (pc PlayerCache) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := yaml.Marshal(pc)
	if err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}
