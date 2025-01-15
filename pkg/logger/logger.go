package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type Config struct {
	Level        string
	IsPrettified bool
}

func New(cfg Config) *slog.Logger {
	var log *slog.Logger

	logLevel := slog.LevelInfo

	switch cfg.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "error":
		logLevel = slog.LevelError
	case "warn":
		logLevel = slog.LevelWarn
	}

	if cfg.IsPrettified {
		log = setupPrettySlog(logLevel)
	} else {
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
		)
	}
	slog.SetDefault(log)
	return log
}

func setupPrettySlog(level slog.Level) *slog.Logger {
	log := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			TimeFormat: time.DateTime,
		}),
	)

	return log
}
