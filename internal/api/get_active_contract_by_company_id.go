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

// GetActiveContractByCompanyId receives query params already bound by the generated Fiber wrapper.
func (s *Server) GetActiveContractByCompanyId(
	c *fiber.Ctx,
	params gen.GetActiveContractByCompanyIdParams,
) error {
	tracer := otel.Tracer("contract-api")
	ctx, span := tracer.Start(c.UserContext(), "GetActiveContractByCompanyIDHandler")
	traceID := span.SpanContext().TraceID().String()
	defer span.End()

	logger := logctx.Logger(ctx).With(
		slog.String("server", "ContractServer"),
		slog.String("handler", "GetActiveContractByCompanyId"),
		slog.String("trace_id", traceID),
	)

	if s.Validator != nil {
		if err := s.Validator.Struct(&params); err != nil {
			logger.Warn("get active contract failed: invalid request", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	span.SetAttributes(attribute.Int64("company_id", params.CompanyId))
	ctx = logctx.WithLogger(ctx, logger)
	logger.Debug("GetActiveContractByCompanyId started")

	resp, err := s.ContractService.GetActiveContractByCompanyID(
		ctx,
		toGetActiveContractByCompanyIDInput(&params),
	)
	if err != nil {
		logger.Error("get active contract failed", slog.Any("error", err))
		return HandleError(c, err)
	}

	out := toContractResponse(resp)
	logger.Debug("get active contract completed", slog.String("company_id", strconv.FormatInt(params.CompanyId, 10)))
	return c.Status(fiber.StatusOK).JSON(out)
}
