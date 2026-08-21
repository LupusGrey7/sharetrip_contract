package api

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetAvailableOfferingByCompanyID(ctx *fiber.Ctx) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "http"),
		slog.String("handler", "GetAvailableOfferingByCompanyID"),
		slog.String("company_id", ctx.Params("companyId")),
		slog.String("service_code", ctx.Params("serviceCode")),
	)
	logger.Debug("GetAvailableOfferingByCompanyID http started")

	var req GetAvailableOfferingByCompanyIDRequest
	if err := ctx.ParamsParser(&req); err != nil {
		logger.Warn("invalid path params", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	}

	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Error("invalid request", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	resp, err := s.CompanyService.GetAvailableOfferingByCompanyID(ctx.UserContext(), toAvailabilityInput(&req))
	if err != nil {
		logger.Error("GetAvailableOfferingByCompanyID failed", slog.Any("error", err))
		return HandleError(ctx, err)
	}

	out := toAvailabilityResponse(resp)
	logger.Debug("GetAvailableOfferingByCompanyID success",
		slog.String("company_id", strconv.Itoa(req.CompanyID)),
		slog.Bool("allowed", out.Allowed),
	)
	return ctx.Status(fiber.StatusOK).JSON(out)
}
