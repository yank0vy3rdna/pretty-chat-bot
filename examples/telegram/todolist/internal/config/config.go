package config

import (
	"github.com/kelseyhightower/envconfig"
)

type BotConfig struct {
	Token        string
	RedisAddress string
}

func GetConfig() (BotConfig, error) {
	var cfg BotConfig
	err := envconfig.Usage("bot", &cfg)
	if err != nil {
		return BotConfig{}, err
	}
	err = envconfig.Process("bot", &cfg)
	if err != nil {
		return BotConfig{}, err
	}
	return cfg, nil
}
