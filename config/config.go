package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Bot  *BotConfig
	Misc *MiscConfig
}

type BotConfig struct {
	Token string
	App   string
	Guild string
}

type MiscConfig struct {
	Cooldown int16
}

func NewConfig() *Config {
	return &Config{
		Bot:  &BotConfig{},
		Misc: &MiscConfig{},
	}
}

func InitConfig(path string) (*Config, error) {
	viper.AddConfigPath(path + "/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := NewConfig()

	cfg.Bot.App = viper.GetString("bot.app")
	cfg.Bot.Guild = viper.GetString("bot.guild")
	cfg.Bot.Token = viper.GetString("bot.token")

	cfg.Misc.Cooldown = int16(viper.GetInt("misc.cooldown"))

	log.Info().Msg("Config successfully loaded")

	return cfg, nil
}
