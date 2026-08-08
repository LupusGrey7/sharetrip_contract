package http

import (
	"log/slog"
	"os"
	"strconv"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetContractByID(ctx *fiber.Ctx) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "http"),
		slog.String("handler", "GetContractByID"),
		slog.String("contract_id", ctx.Params("contractId")),
	)
	logger.Debug("GetContractByID http started")

	contractID := ctx.Params("contractId")
	if contractID == "" {
		logger.Warn("get contract by id failed: invalid request", slog.String("error", ErrInvalidIDParamFormat.Error()))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	}

	contractIDInt, err := strconv.Atoi(contractID)
	if err != nil {
		logger.Error("get contract by id failed: invalid request", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	}

	req := model.GetContractByIDRequest{
		ContractID: contractIDInt,
	}

	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Error("get contract by id failed: invalid request", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	resp, err := s.ContractService.GetContractByID(ctx.UserContext(), &req)
	if err != nil {
		logger.Error("get contract by id failed", slog.Any("error", err))
		return HandleError(ctx, err)
	}

	logger.Debug("get contract by id completed", slog.String("contract_id", contractID))
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
