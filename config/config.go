package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Bot      *BotConfig
	Misc     *MiscConfig
	DBConfig *DBConfig
}

type BotConfig struct {
	Token string
	App   string
	Guild string
}

type MiscConfig struct {
	UpdateTime      int64
	DisconnectTimer int64
	ClearingTimer   int64
	Cooldown        int16
	CacheSize       int16
	QueueSize       int16
	Folder          string
}

type DBConfig struct {
	Url      string
	User     string
	Password string
	Port     int
}

func NewConfig() *Config {
	return &Config{
		Bot:      &BotConfig{},
		Misc:     &MiscConfig{},
		DBConfig: &DBConfig{},
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

	cfg.DBConfig.Url = viper.GetString("db.url")
	cfg.DBConfig.User = viper.GetString("db.user")
	cfg.DBConfig.Password = viper.GetString("db.password")
	cfg.DBConfig.Port = viper.GetInt("db.port")

	cfg.Misc.Cooldown = int16(viper.GetInt("misc.cooldown"))
	cfg.Misc.CacheSize = int16(viper.GetInt("misc.cacheSize"))
	cfg.Misc.QueueSize = int16(viper.GetInt("misc.queueSize"))
	cfg.Misc.UpdateTime = int64(viper.GetInt("misc.updateTime"))
	cfg.Misc.DisconnectTimer = int64(viper.GetInt("misc.disconnectTimer"))
	cfg.Misc.ClearingTimer = int64(viper.GetInt("misc.clearingTimer"))
	cfg.Misc.Folder = viper.GetString("misc.folder")

	log.Info().Msg("Config successfully loaded")

	return cfg, nil
}
