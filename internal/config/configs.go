package config

import (
	"fmt"
	"time"
)

type AppConfig struct {
	Env  string `env:"ENV" mapstructure:"env"`
	Name string `env:"NAME" mapstructure:"name"`
}

func (a *AppConfig) IsDev() bool {
	return a.Env == "local" || a.Env == "dev"
}

type LoggerConfig struct {
	Level        string `env:"LEVEL" mapstructure:"level"`
	IsPrettified bool   `env:"IS_PRETTIFIED" mapstructure:"is_prettified"`
}

// Adapters ============================

type AdaptersConfig struct {
	Redis AdapterRedisConfig `envPrefix:"REDIS_" mapstructure:"redis"`
}

type BaseRetriesOptionsConfig struct {
	Attempts int           `env:"ATTEMPTS" mapstructure:"attempts"`
	Interval time.Duration `env:"INTERVAL" mapstructure:"interval"`
}

type AdapterRedisConfig struct {
	Host     string                   `env:"HOST" mapstructure:"host"`
	Port     int                      `env:"PORT" mapstructure:"port"`
	Password string                   `env:"PASSWORD" mapstructure:"password"`
	DB       int                      `env:"DB" mapstructure:"db"`
	Pool     AdapterRedisPoolConfig   `envPrefix:"POOL_" mapstructure:"pool"`
	Retries  BaseRetriesOptionsConfig `envPrefix:"RETRIES_" mapstructure:"retries"`
}

type AdapterRedisPoolConfig struct {
	MinSize     int           `env:"MIN_SIZE" mapstructure:"min_size"`
	MaxSize     int           `env:"MAX_SIZE" mapstructure:"max_size"`
	MinIdleSize int           `env:"MIN_IDLE_SIZE" mapstructure:"min_idle_size"`
	MaxIdleSize int           `env:"MAX_IDLE_SIZE" mapstructure:"max_idle_size"`
	MaxIdleTime time.Duration `env:"MAX_IDLE_TIME" mapstructure:"max_idle_time"`
}

func (c *AdapterRedisConfig) Addr() string {
	return fmt.Sprintf(
		"%s:%d",
		c.Host,
		c.Port,
	)
}

// Transports ==========================

type TransportsConfig struct {
	HTTP BaseHostPortConfig `envPrefix:"HTTP_" mapstructure:"http"`
}

type BaseHostPortConfig struct {
	Host string `env:"HOST" mapstructure:"host"`
	Port int    `env:"PORT" mapstructure:"port"`
}

func (c *BaseHostPortConfig) Addr() string {
	return fmt.Sprintf(
		"%s:%d",
		c.Host,
		c.Port,
	)
}
