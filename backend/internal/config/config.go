package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	PgDSN         string `mapstructure:"PG_DSN"`
	MigratePath   string `mapstructure:"MIGRATE_PATH"`
}

func New(path string) (Config, error) {
	var cfg Config

	viper.SetConfigName(".env")
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
