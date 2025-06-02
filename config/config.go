package config

import (
	"github.com/spf13/viper"
)

var GlobalConfig Config

type Config struct {
	DBdriver               string `mapstructure:"DB_DRIVER"`
	DBSource               string `mapstructure:"DB_SOURCE"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	CampaignServiceAddress string `mapstructure:"CAMPAIGN_SERVICE_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	err = viper.Unmarshal(&GlobalConfig)
	return
}
