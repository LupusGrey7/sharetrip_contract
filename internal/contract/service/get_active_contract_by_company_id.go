package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/observability/logctx"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func (s *ContractService) GetActiveContractByCompanyID(
	ctx context.Context,
	input *domain.GetActiveContractByCompanyIDInput,
) (res *domain.ContractOutput, err error) {
	ctxSpc, span := otel.Tracer("ContractService").Start(ctx, "ContractService.GetActiveContractByCompanyID")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	logger := logctx.Logger(ctxSpc).With(
		slog.String("service", "ContractService"),
		slog.String("operation", "GetActiveContractByCompanyID"),
		slog.String("company_id", strconv.Itoa(input.CompanyID)),
	)
	logger.Debug("get active contract started")

	if s.pool == nil {
		res, err = s.useCase.GetActiveContractByCompanyID(ctxSpc, nil, s.repo, input)
		if err != nil {
			return nil, fmt.Errorf("GetActiveContractByCompanyID: %w", err)
		}
		logger.Debug("get active contract completed", slog.String("contract_id", res.ID.String()))
		return res, nil
	}

	txCtx, txSpan := otel.Tracer("database").Start(ctxSpc, "DB.Transaction")
	defer txSpan.End()

	res, err = tx(txCtx, s.pool, func(pgTx pgx.Tx) (*domain.ContractOutput, error) {
		return s.useCase.GetActiveContractByCompanyID(txCtx, pgTx, s.repo, input)
	})
	if err != nil {
		return nil, fmt.Errorf("GetActiveContractByCompanyID: %w", err)
	}

	logger.Debug("get active contract completed", slog.String("contract_id", res.ID.String()))
	return res, nil
}
