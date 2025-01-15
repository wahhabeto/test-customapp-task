package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"

	"github.com/FGP-tech/common/pkg/utils/retryable"
)

type RetriesConfig struct {
	Attempts int
	Interval time.Duration
}

type Config struct {
	Addr            string
	Password        string
	DB              int
	PoolMinSize     int
	PoolMaxSize     int
	PoolMinIdleSize int
	PoolMaxIdleSize int
	PoolMaxIdleTime time.Duration

	Retries RetriesConfig
}

func NewClient(log *slog.Logger, cfg Config) (*redis.Client, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:            cfg.Addr,
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolMinSize,
		MaxActiveConns:  cfg.PoolMaxSize,
		MinIdleConns:    cfg.PoolMinIdleSize,
		MaxIdleConns:    cfg.PoolMaxIdleSize,
		ConnMaxIdleTime: cfg.PoolMaxIdleTime,
	})

	err := retryable.DoWithRetry(func() error {
		_, err := client.Ping(ctx).Result()
		return err
	}, retryable.WithRetryOptions{
		Attempts: cfg.Retries.Attempts,
		Interval: cfg.Retries.Interval,
		Logger:   log,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info("Connected")

	return client, nil
}
