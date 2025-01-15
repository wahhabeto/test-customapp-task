package di

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/wahhabeto/test-customapp-task/internal/app"
	"github.com/wahhabeto/test-customapp-task/internal/config"
	"github.com/wahhabeto/test-customapp-task/internal/transport/http"
	"github.com/wahhabeto/test-customapp-task/pkg/locker"
	"go.uber.org/dig"
	"reflect"
	"time"
)

type DIGroupFunc func(c *dig.Container, groupName string)

func MustCreateDiContainer(rtp float64) *dig.Container {
	_ = godotenv.Load(".env")

	c := dig.New()

	constructors := []interface{}{
		// Конфигурации ===========
		config.NewConfig,

		// Locker ====================
		locker.NewLocker,

		// Logger ===================
		newLogger,

		// Adapters =================
		resty.New,
		newRedis,

		// engine
		func(cfg *config.Config) (app.Engine, error) {
			// Передаем rtp из аргумента в функцию NewEngine
			return app.NewEngine(rtp)
		},

		// Servers ==================
		http.NewServer,

		// ==========================
		NewApp,
	}

	// Регистрация зависимостей
	for _, constructor := range constructors {
		constructorType := reflect.TypeOf(constructor).String()

		if err := c.Provide(constructor); err != nil {
			log.Errorf("Failed to register constructor %s: %v", constructorType, err)
			panic(fmt.Errorf("failed to provide dependency: %w", err))
		}
	}

	_ = c.Decorate(func(r *resty.Client) *resty.Client {
		return r.SetTimeout(time.Minute * 3).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	})

	// Graceful shutdown
	closers := []interface{}{
		closeHTTP,
	}

	for _, closeFunc := range closers {
		_ = c.Invoke(closeFunc)
	}

	return c
}
