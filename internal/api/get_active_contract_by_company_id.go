package api

import (
	"log/slog"
	"strconv"

	"job4j/sharetrip-contract/internal/observability/logctx"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Server) GetActiveContractByCompanyID(c *fiber.Ctx) error {
	tracer := otel.Tracer("contract-api")
	ctx, span := tracer.Start(c.UserContext(), "GetActiveContractByCompanyIDHandler")
	traceID := span.SpanContext().TraceID().String()
	defer span.End()

	logger := logctx.Logger(ctx).With(
		slog.String("server", "ContractServer"),
		slog.String("handler", "GetActiveContractByCompanyID"),
		slog.String("trace_id", traceID),
	)

	var req GetActiveContractByCompanyIDRequest
	if err := c.QueryParser(&req); err != nil {
		logger.Warn("get active contract failed: invalid query", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	}
	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Warn("get active contract failed: invalid request", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	span.SetAttributes(attribute.Int("company_id", req.CompanyID))
	ctx = logctx.WithLogger(ctx, logger)
	logger.Debug("GetActiveContractByCompanyID started")

	resp, err := s.ContractService.GetActiveContractByCompanyID(
		ctx,
		toGetActiveContractByCompanyIDInput(&req),
	)
	if err != nil {
		logger.Error("get active contract failed", slog.Any("error", err))
		return HandleError(c, err)
	}

	out := toContractResponse(resp)
	logger.Debug("get active contract completed", slog.String("company_id", strconv.Itoa(req.CompanyID)))
	return c.Status(fiber.StatusOK).JSON(out)
}
