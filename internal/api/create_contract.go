package api

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"job4j/sharetrip-contract/internal/observability/logctx"
)

func (s *Server) CreateContract(c *fiber.Ctx) error {
	tracer := otel.Tracer("trip-contract-api")
	ctx, span := tracer.Start(c.UserContext(), "CreateContractHandler")
	traceID := span.SpanContext().TraceID().String()
	c.Set("X-Request-ID", traceID)
	defer span.End()

	logger := logctx.Logger(ctx).With(
		slog.String("server", "TripContractServer"),
		slog.String("handler", "CreateContract"),
		slog.String("trace_id", traceID),
	)

	var req CreateContractRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warn("CreateContract parse failed", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidRequest.Error())
	}

	if s.Validator != nil {
		if err := s.Validator.Struct(&req); err != nil {
			logger.Warn("CreateContract validate failed", slog.Any("error", err))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	ctx = logctx.WithLogger(ctx, logger)
	logger.Debug("CreateContract started")

	resp, err := s.ContractService.CreateContract(ctx, toCreateContractInput(&req))
	if err != nil {
		logger.Error("CreateContract failed", slog.Any("error", err))
		return HandleError(c, err)
	}

	out := toContractResponse(resp)
	logger.Debug("CreateContract completed", slog.Int("contract_id", out.ID))
	return c.Status(fiber.StatusCreated).JSON(out)
}
