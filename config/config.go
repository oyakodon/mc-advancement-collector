package config

const (
	CONFIG_PATH          = "./config/config.yml"
	ADVANCEMENTLIST_PATH = "./config/advancementlist.yml"
	PLAYERCACHE_PATH     = "./config/playercache.yml"
)

type Config struct {
	AppConfig       *AppConfig
	AdvancementList *AdvancementList
	PlayerCache     *PlayerCache
}

func LoadConfig() (*Config, error) {
	conf, err := LoadAppConfig(CONFIG_PATH)
	if err != nil {
		return nil, err
	}

	advancementlist, err := LoadAdvancementList(ADVANCEMENTLIST_PATH)
	if err != nil {
		return nil, err
	}

	playercache, err := LoadPlayerCache(PLAYERCACHE_PATH)
	if err != nil {
		return nil, err
	}

	return &Config{
		AppConfig:       conf,
		AdvancementList: advancementlist,
		PlayerCache:     playercache,
	}, nil
}
