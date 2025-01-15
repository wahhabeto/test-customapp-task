package locker

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
	"time"
)

type Locker interface {
	NewLock(key string, ttlSec int) *Lock
}

type locker struct {
	client *redis.Client
	log    *slog.Logger
}

func NewLocker(client *redis.Client, log *slog.Logger) Locker {
	return &locker{
		client: client,
		log:    log,
	}
}

func (l *locker) NewLock(key string, ttlSec int) *Lock {
	return &Lock{
		client:   l.client,
		key:      key,
		value:    uuid.New().String(),
		duration: time.Duration(ttlSec) * time.Second,
	}
}

// ---------------------------- LOCK ----------------------------------

type Lock struct {
	client   *redis.Client
	key      string
	value    string
	duration time.Duration
}

func (l *Lock) Acquire(ctx context.Context) (bool, error) {

	duration := strconv.FormatInt(l.duration.Milliseconds(), 10)

	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
			return "OK"
		else
			return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
		end
	`

	cmd, err := l.client.Eval(ctx, script, []string{l.key}, []string{l.value, duration}).Result()
	if err != nil {
		slog.Warn("Lock failed", "cmd", cmd, "err", err.Error())
		return false, err
	}
	slog.Debug("Lock complete", "cmd", cmd)
	if cmd == "OK" {
		return true, nil
	}

	return false, nil
}

func (l *Lock) Release(ctx context.Context) error {
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	_, err := l.client.Eval(ctx, script, []string{l.key}, []string{l.value}).Result()
	return err
}
