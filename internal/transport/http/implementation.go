package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/wahhabeto/test-customapp-task/internal/transport/http/api"
)

func (s *server) HealthCheck(c *fiber.Ctx) error {
	s.log.Debug("Root handler invoked")
	return c.SendStatus(fiber.StatusOK)
}

func (s *server) GetGeneratedNumber(c *fiber.Ctx) error {
	//todo:import lock config from configs problm
	lock := s.locker.NewLock("key", 10) // todo:add config for lock timeout
	isLocked, lockErr := lock.Acquire(c.Context())
	if !isLocked {
		return badRequest(c)
	}
	if lockErr != nil {
		return badRequest(c)
	}
	defer func(ctx context.Context) {
		if err := lock.Release(ctx); err != nil {
		}
	}(c.Context())

	multiplier, err := s.eng.GenerateMultiplier()
	if err != nil {
		s.log.Error("Error generating multiplier: %v", err)
		return badRequest(c)
	}

	return c.JSON(api.GetResponse{
		Result: multiplier,
	})
}

func badRequest(c *fiber.Ctx) error {
	_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "Bad Request",
		"id":    c.Get("RequestID", ""),
	})
	return nil
}
