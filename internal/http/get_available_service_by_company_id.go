package http

import (
	"log/slog"
	"os"
	"strconv"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetAvailableServiceByCompanyID(ctx *fiber.Ctx) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "http"),
		slog.String("handler", "GetAvailableServiceByCompanyID"),
		slog.String("company_id", ctx.Params("companyId")),
		slog.String("service_code", ctx.Params("serviceCode")),
	)
	logger.Debug("GetAvailableServiceByCompanyID http started")

	companyID := ctx.Params("companyId")
	if companyID == "" {
		logger.Warn("get available service by company id failed: invalid request", slog.String("error", ErrInvalidIDParamFormat.Error()))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	}

	companyIDInt, err := strconv.Atoi(companyID)
	if err != nil {
		logger.Error("get available service by company id failed: invalid request", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	}

	serviceCode := ctx.Params("serviceCode")
	if serviceCode == "" {
		logger.Warn("get available service by company id failed: invalid request", slog.String("error", ErrInvalidRequest.Error()))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidRequest.Error())
	}

	modelRequest := &model.GetAvailableServiceByCompanyIDRequest{
		CompanyID:   companyIDInt,
		ServiceCode: model.ServiceCodeType(serviceCode),
	}

	if s.Validator != nil {
		if err := s.Validator.Struct(modelRequest); err != nil {
			logger.Error("get available service by company id failed: invalid request", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	resp, err := s.service.GetAvailableServiceByCompanyID(ctx.UserContext(), modelRequest)
	if err != nil {
		logger.Error("get available service by company id failed", slog.Any("error", err))
		return HandleError(ctx, err)
	}

	out := toAvailabilityResponse(resp)
	logger.Debug("get available service by company id success", slog.Any("response", out))
	return ctx.Status(fiber.StatusOK).JSON(out)
}
