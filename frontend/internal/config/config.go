package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress     string        `mapstructure:"SERVER_ADDRESS"`
	APIurl            string        `mapstructure:"API_URL"`
	APIRequestTimeout time.Duration `mapstructure:"API_REQUEST_TIMEOUT"`
}

func New(path string) (Config, error) {
	var cfg Config

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.SetDefault("API_REQUEST_TIMEOUT", "1s")
	viper.SetDefault("ITERATION_DELAY", "5s")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
