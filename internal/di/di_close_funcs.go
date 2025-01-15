package di

import (
	"github.com/wahhabeto/test-customapp-task/internal/transport/http"
	"github.com/wahhabeto/test-customapp-task/pkg/closer"
	"log/slog"
)

// Common ============================

func closeHTTP(s http.Server, log *slog.Logger) {
	closer.Add(func() error {
		log.Debug("[Shutdown]", slog.String("component", "HTTP"))
		return s.Close()
	})
}
