package http

import (
	"log/slog"
	"os"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) UpsertContractServices(ctx *fiber.Ctx) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "http"),
		slog.String("handler", "UpsertContractServices"),
	)
	logger.Debug("UpsertContractServices started")

	var req model.UpsertContractServicesRequest
	if err := ctx.BodyParser(&req); err != nil {
		logger.Warn("UpsertContractServices parse failed", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidRequest.Error())
	}

	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Warn("UpsertContractServices validate failed", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	resp, err := s.OfferingService.UpsertContractServices(ctx.UserContext(), &req)
	if err != nil {
		logger.Error("UpsertContractServices failed", slog.Any("error", err))
		return HandleError(ctx, err)
	}

	logger.Debug("UpsertContractServices completed", slog.Int("contract_id", resp.ContractID))
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
