package api

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetContractByID(ctx *fiber.Ctx) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "http"),
		slog.String("handler", "GetContractByID"),
		slog.String("contract_id", ctx.Params("contractId")),
	)
	logger.Debug("GetContractByID http started")

	var req GetContractByIDRequest
	if err := ctx.ParamsParser(&req); err != nil {
		logger.Warn("get contract by id failed: invalid path params", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	}
	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Error("get contract by id failed: invalid request", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	resp, err := s.ContractService.GetContractByID(ctx.UserContext(), toGetContractByIDInput(&req))
	if err != nil {
		logger.Error("get contract by id failed", slog.Any("error", err))
		return HandleError(ctx, err)
	}

	out := toContractResponse(resp)
	logger.Debug("get contract by id completed", slog.String("contract_id", strconv.Itoa(req.ContractID)))
	return ctx.Status(fiber.StatusOK).JSON(out)
}
