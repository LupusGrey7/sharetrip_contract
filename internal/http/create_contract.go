package http

import (
	"log/slog"
	"os"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) CreateContract(ctx *fiber.Ctx) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "http"),
		slog.String("handler", "CreateContract"),
	)
	logger.Debug("CreateContract started")

	var req model.CreateContractRequest
	if err := ctx.BodyParser(&req); err != nil {
		logger.Warn("CreateContract parse failed", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidRequest.Error())
	}

	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Warn("CreateContract validate failed", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	resp, err := s.ContractService.CreateContract(ctx.UserContext(), &req)
	if err != nil {
		logger.Error("CreateContract failed", slog.Any("error", err))
		return HandleError(ctx, err)
	}

	logger.Debug("CreateContract completed", slog.Int("contract_id", resp.ID))
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}
