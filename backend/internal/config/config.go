package config

import "github.com/spf13/viper"

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	PgDSN         string `mapstructure:"PG_DSN"`
	MigratePath   string `mapstructure:"MIGRATE_PATH"`
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
