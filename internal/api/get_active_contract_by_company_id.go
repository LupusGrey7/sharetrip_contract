package api

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetActiveContractByCompanyID(ctx *fiber.Ctx) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "http"),
		slog.String("handler", "GetActiveContractByCompanyID"),
	)
	logger.Debug("GetActiveContractByCompanyID http started")

	var req GetActiveContractByCompanyIDRequest
	if err := ctx.QueryParser(&req); err != nil {
		logger.Warn("get active contract failed: invalid query", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	}
	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Error("get active contract failed: invalid request", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	resp, err := s.ContractService.GetActiveContractByCompanyID(
		ctx.UserContext(),
		toGetActiveContractByCompanyIDInput(&req),
	)
	if err != nil {
		logger.Error("get active contract failed", slog.Any("error", err))
		return HandleError(ctx, err)
	}

	out := toContractResponse(resp)
	logger.Debug("get active contract completed", slog.String("company_id", strconv.Itoa(req.CompanyID)))
	return ctx.Status(fiber.StatusOK).JSON(out)
}
