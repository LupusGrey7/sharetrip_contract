package api

import (
	"log/slog"

	"job4j/sharetrip-contract/gen"
	"job4j/sharetrip-contract/internal/observability/logctx"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func (s *Server) CreateContract(c *fiber.Ctx) error {
	tracer := otel.Tracer("contract-api")
	ctx, span := tracer.Start(c.UserContext(), "CreateContractHandler")
	defer span.End()

	logger := logctx.Logger(ctx).With(
		slog.String("server", "ContractServer"),
		slog.String("handler", "CreateContract"),
	)

	// Wire type from gen (oapi-codegen). Logic same as before: parse → Validator → service.
	var req gen.CreateContractRequest
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

	out := toCreateContractResponse(resp)
	logger.Debug("CreateContract completed completed", slog.String("contract_id", out.Contract.Id.String()))
	return c.Status(fiber.StatusCreated).JSON(out)
}
