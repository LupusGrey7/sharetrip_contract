package api

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) Healthcheck(ctx *fiber.Ctx) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "http"),
		slog.String("handler", "Healthcheck"),
	)
	logger.Debug("healthcheck started")

	resp, err := s.HealthcheckService.GetHealthcheckInfo(ctx.UserContext())
	if err != nil {
		logger.Error("healthcheck failed", slog.Any("error", err))
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(HealthcheckResponse{
			Status:  "error",
			Message: "database unavailable",
		})
	}

	logger.Debug("healthcheck completed")
	return ctx.Status(fiber.StatusOK).JSON(toHealthcheckResponse(resp))
}
