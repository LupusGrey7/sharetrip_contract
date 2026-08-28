package api

import (
	"log/slog"
	"strconv"

	"job4j/sharetrip-contract/internal/observability/logctx"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Server) GetAvailableOfferingByCompanyID(c *fiber.Ctx) error {
	tracer := otel.Tracer("contract-api")
	ctx, span := tracer.Start(c.UserContext(), "GetAvailableOfferingByCompanyIDHandler")
	traceID := span.SpanContext().TraceID().String()
	defer span.End()

	logger := logctx.Logger(ctx).With(
		slog.String("server", "ContractServer"),
		slog.String("handler", "GetAvailableOfferingByCompanyID"),
		slog.String("trace_id", traceID),
		slog.String("company_id", c.Params("companyId")),
		slog.String("service_code", c.Params("serviceCode")),
	)

	var req GetAvailableOfferingByCompanyIDRequest
	if err := c.ParamsParser(&req); err != nil {
		logger.Warn("invalid path params", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	}

	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Warn("invalid request", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	span.SetAttributes(
		attribute.Int("company_id", req.CompanyID),
		attribute.String("service_code", string(req.ServiceCode)),
	)
	ctx = logctx.WithLogger(ctx, logger)
	logger.Debug("GetAvailableOfferingByCompanyID started")

	resp, err := s.CompanyService.GetAvailableOfferingByCompanyID(ctx, toAvailabilityInput(&req))
	if err != nil {
		logger.Error("GetAvailableOfferingByCompanyID failed", slog.Any("error", err))
		return HandleError(c, err)
	}

	out := toAvailabilityResponse(resp)
	logger.Debug("GetAvailableOfferingByCompanyID success",
		slog.String("company_id", strconv.Itoa(req.CompanyID)),
		slog.Bool("allowed", out.Allowed),
	)
	return c.Status(fiber.StatusOK).JSON(out)
}
