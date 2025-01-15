package di

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wahhabeto/test-customapp-task/internal/config"
	oredis "github.com/wahhabeto/test-customapp-task/pkg/redis"
	"log/slog"

	"github.com/wahhabeto/test-customapp-task/pkg/logger"
)

// Common ============================

func newLogger(cfg *config.Config) *slog.Logger {
	return logger.New(logger.Config{
		Level:        cfg.Logger.Level,
		IsPrettified: cfg.Logger.IsPrettified,
	})
}

func newRedis(log *slog.Logger, cfg *config.Config) *redis.Client {
	client, err := oredis.NewClient(
		log.With(
			slog.String("component", "redis"),
		),
		oredis.Config{
			Addr:            cfg.Adapters.Redis.Addr(),
			Password:        cfg.Adapters.Redis.Password,
			DB:              cfg.Adapters.Redis.DB,
			PoolMinSize:     cfg.Adapters.Redis.Pool.MinSize,
			PoolMaxSize:     cfg.Adapters.Redis.Pool.MaxSize,
			PoolMinIdleSize: cfg.Adapters.Redis.Pool.MinIdleSize,
			PoolMaxIdleSize: cfg.Adapters.Redis.Pool.MaxIdleSize,
			PoolMaxIdleTime: cfg.Adapters.Redis.Pool.MaxIdleTime,
			Retries: oredis.RetriesConfig{
				Attempts: cfg.Adapters.Redis.Retries.Attempts,
				Interval: cfg.Adapters.Redis.Retries.Interval,
			},
		})
	if err != nil {
		panic(fmt.Errorf("failed to create redis client: %w", err))
	}
	return client
}
