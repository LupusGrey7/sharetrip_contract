package service

import (
	"context"
	"fmt"
	"log/slog"
0
	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"job4j/sharetrip-contract/internal/observability/logctx"
)

func (s *ContractService) CreateContract(
	ctx context.Context,
	input *domain.CreateContractInput,
) (res *domain.ContractOutput, err error) {
	// 1. Integration with Jaeger: create a child span for this layer (ctx now contains the ID of this span)
	ctxSpc, span := otel.Tracer("TripService").Start(ctx, "TripService.CreateTripDraft")

	defer func() {
		// Now 'err' is taken from the return of the function. If an error occurred below in the stack, err != nil
		if err != nil {
			span.RecordError(err)                    // Log error in Jaeger
			span.SetStatus(codes.Error, err.Error()) // Color the span in Jaeger in red color
		}

		span.End() // Span always ends in the end. Jaeger will measure the time between Start and End!
	}()

	//getting custom logger context
	logger := logctx.Logger(ctxSpc).With(
		slog.String("service", "TripContractService"),
		slog.String("operation", "CreateTripDraft"),
		slog.String("company_id", req.CompanyID.String()),
	)
	logger.Debug("create trip draft started")

	// 4. Open a database transaction
	// Mandatory to create a sub-span for the transaction to measure its clean duration
	txCtx, txSpan := otel.Tracer("database").Start(ctxSpc, "DB.Transaction")
	defer txSpan.End()

	res, err = tx(txCtx, s.pool, func(pgTx pgx.Tx) (*domain.ContractOutput, error) {
		txLogger := logger.With(slog.String("layer", "transaction"))
		txLogger.Debug("transaction create contract execution started")

		resp, err := s.useCase.CreateContract(ctx, pgTx, s.repo, input)
		if err != nil {
			txLogger.Error("create trip contract usecase failed", slog.Any("error", err))
			return nil, fmt.Errorf("usecase.CreateContract: %w", err)
		}

		txLogger.Debug("transaction create contract completed", slog.String("contract_id", resp.ID.String()))
		return resp, nil
	})
	if err != nil {
		return nil, fmt.Errorf("CreateContract: %w", err)
	}

	logger.Debug("create contract completed", slog.String("trip_id", res.ID.String()))
	return res, nil
}
