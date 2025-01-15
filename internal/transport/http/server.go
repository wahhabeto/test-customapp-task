package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
	"github.com/wahhabeto/test-customapp-task/internal/app"
	"github.com/wahhabeto/test-customapp-task/internal/config"
	"github.com/wahhabeto/test-customapp-task/pkg/locker"
	"go.uber.org/dig"
	"log/slog"
	"time"
)

type Server interface {
	Run(ctx context.Context) error
	Close() error
}

type server struct {
	app *fiber.App

	eng    app.Engine
	locker locker.Locker
	log    *slog.Logger
	conf   *config.Config
}

type serverParams struct {
	dig.In

	Eng  app.Engine
	Log  *slog.Logger
	Conf *config.Config

	Locker locker.Locker
}

func NewServer(params serverParams) Server {
	log := params.Log.With(
		slog.String("component", "http"),
	)

	app := fiber.New(fiber.Config{
		ServerHeader: params.Conf.App.Name,
		ReadTimeout:  time.Duration(5) * time.Second,
		WriteTimeout: time.Duration(5) * time.Second,
		IdleTimeout:  time.Duration(60) * time.Second,
	})

	app.Use(slogfiber.New(log))

	app.Use(recover.New())

	//app.Use(swagger.New(swagger.Config{
	//	BasePath: "/api/v1/",
	//	FilePath: "./docs/swagger.json",
	//	Path:     "swagger",
	//	Title:    "Swagger API Docs",
	//}))

	return &server{
		app: app,

		eng: params.Eng,

		log: log,

		conf: params.Conf,

		locker: params.Locker,
	}
}

func (s *server) Run(ctx context.Context) error {
	address := s.conf.Transports.HTTP.Addr()

	if err := s.MountRouter(); err != nil {
		return err
	}

	s.log.Info("Starting HTTP server", "address", address)

	return s.app.Listen(address)
}

func (s *server) Close() error {
	return s.app.Shutdown()
}
