package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APIurl            string        `mapstructure:"API_URL"`
	APIRequestTimeout time.Duration `mapstructure:"API_REQUEST_TIMEOUT"`
	IterationDelay    time.Duration `mapstructure:"ITERATION_DELAY"`
	Count             int           `mapstructure:"COUNT"`
	Interval          time.Duration `mapstructure:"INTERVAL"`
	Timeout           time.Duration `mapstructure:"TIMEOUT"`
}

func New(path string) (Config, error) {
	var cfg Config

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.SetDefault("API_REQUEST_TIMEOUT", "1s")
	viper.SetDefault("ITERATION_DELAY", "5s")
	viper.SetDefault("COUNT", 5)
	viper.SetDefault("INTERVAL", "100ms")
	viper.SetDefault("TIMEOUT", "1s")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
