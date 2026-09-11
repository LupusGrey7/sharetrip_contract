package api

import (
	"log/slog"
	"strconv"

	"job4j/sharetrip-contract/gen"
	"job4j/sharetrip-contract/internal/observability/logctx"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CheckServiceAvailabilityForCompany receives path params already bound by the generated Fiber wrapper.
func (s *Server) CheckServiceAvailabilityForCompany(
	c *fiber.Ctx,
	companyId int64,
	serviceCode gen.CheckServiceAvailabilityForCompanyParamsServiceCode,
) error {
	tracer := otel.Tracer("contract-api")
	ctx, span := tracer.Start(c.UserContext(), "GetAvailableOfferingByCompanyIDHandler")
	traceID := span.SpanContext().TraceID().String()
	defer span.End()

	logger := logctx.Logger(ctx).With(
		slog.String("server", "ContractServer"),
		slog.String("handler", "CheckServiceAvailabilityForCompany"),
		slog.String("trace_id", traceID),
		slog.Int64("company_id", companyId),
		slog.String("service_code", string(serviceCode)),
	)

	// Wrapper has already converted path strings to typed values.
	// This request object keeps the existing min/oneof Validator rules.
	req := GetAvailableOfferingByCompanyIDRequest{
		CompanyID:   int(companyId),
		ServiceCode: string(serviceCode),
	}
	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Warn("invalid request", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	span.SetAttributes(
		attribute.Int("company_id", req.CompanyID),
		attribute.String("service_code", req.ServiceCode),
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
