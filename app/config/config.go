package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PRIVATE_KEY_ID string `mapstructure:"PRIVATE_KEY_ID"`
	PRIVATE_KEY    string `mapstructure:"PRIVATE_KEY"`
	CLIENT_ID      string `mapstructure:"CLIENT_ID"`
	CLIENT_EMAIL   string `mapstructure:"CLIENT_EMAIL"`
	CLIENT_URL     string `mapstructure:"CLIENT_URL"`
	APP_ID         string `mapstructure:"APP_ID"`
	MESSAGING_SENDER_ID string `mapstructure:"MESSAGING_SENDER_ID"`
	PROJECT_ID string `mapstructure:"PROJECT_ID"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
