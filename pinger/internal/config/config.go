package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APIurl         string        `mapstructure:"API_URL"`
	RequestTimeout time.Duration `mapstructure:"REQUEST_TIMEOUT"`
}

func New(path string) (Config, error) {
	var cfg Config

	viper.SetConfigName(".app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
