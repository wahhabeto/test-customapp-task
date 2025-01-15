package di

import (
	"context"
	"log/slog"

	"github.com/wahhabeto/test-customapp-task/internal/transport/http"
	"go.uber.org/dig"
)

type App struct {
	log *slog.Logger

	httpServer http.Server
}

type params struct {
	dig.In

	Log *slog.Logger

	HttpServer http.Server
}

func NewApp(p params) *App {
	return &App{
		log: p.Log.With(
			slog.String("component", "app"),
		),

		httpServer: p.HttpServer,
	}
}

func (a *App) Run(ctx context.Context) error {
	ch := make(chan error)

	go func(ctx context.Context, ch chan error) {
		ch <- a.httpServer.Run(ctx)
	}(ctx, ch)

	return <-ch
}
