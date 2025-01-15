package config

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `envPrefix:"APP_" mapstructure:"app"`
	Logger LoggerConfig `envPrefix:"LOG_" mapstructure:"logging"`

	Adapters AdaptersConfig `envPrefix:"ADAPTERS_" mapstructure:"adapters"`

	Transports TransportsConfig `envPrefix:"TRANSPORTS_" mapstructure:"transports"`
}

var configPath string

func NewConfig() (*Config, error) {

	if configPath == "" {
		flag.StringVar(&configPath, "config", "./config.base.yml", "Config file path")
		flag.Parse()
	}
	slog.Debug(
		"Config path",
		slog.String("path", configPath),
	)

	var config Config

	viper.SetConfigFile("./config/config.base.yml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base config: %w", err)
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		slog.Debug(
			"Config file not found",
			slog.String("path", configPath),
		)
	} else {
		viper.SetConfigFile(configPath)

		if err := viper.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	opts := env.Options{RequiredIfNoDef: false}

	if err := env.ParseWithOptions(&config, opts); err != nil {
		return nil, err
	}

	return &config, nil
}
