package retryable

import (
	"fmt"
	"log/slog"
	"time"
)

type WithRetryOptions struct {
	Attempts int
	Interval time.Duration
	Logger   *slog.Logger
}

func DoWithRetry(f func() error, opts ...WithRetryOptions) error {
	var options WithRetryOptions

	if len(opts) > 0 {
		options = opts[0]
	} else {
		options = WithRetryOptions{
			Attempts: 3,
			Interval: 5 * time.Second,
		}

	}
	log := options.Logger
	if log == nil {
		log = slog.Default()
	}

	var err error
	for attempt := 1; attempt <= options.Attempts; attempt++ {
		err = f()
		if err == nil {
			return nil
		}
		log.Info("retrying",
			slog.Int("attempt", attempt),
			slog.String("err", err.Error()),
		)
		time.Sleep(options.Interval)
	}
	return fmt.Errorf("failed after %d attempts: %w", options.Attempts, err)
}
