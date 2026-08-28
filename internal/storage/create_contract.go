package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/observability/logctx"
)

const createContract = `
INSERT INTO contract_management.contracts (
	contract_number, company_id, status_id, start_date, end_date
) VALUES ($1, $2, $3, $4, $5)
RETURNING id, contract_number, company_id, status_id, start_date, end_date, created_at, updated_at`

func (r *ContractRepository) CreateContractTx(
	ctx context.Context,
	tx pgx.Tx,
	t *domain.ContractEntity,
) (*domain.ContractEntity, error) {
	//tracing Jaeger
	tracer := otel.Tracer("ContractRepository")
	ctxSpc, span := tracer.Start(ctx, "ContractRepository.CreateContractTx")

	defer func() {
		// 	rows.Close() // process rows sql FIXME
		span.End() // Span always ends in the end. Jaeger will measure the time between Start and End!
	}()

	//getting custom logger context
	logger := logctx.Logger(ctxSpc).With(
		slog.String("layer", "repository"),
		slog.String("repository", "CreateContractTx"),
		slog.Int("company_id", t.CompanyID),
	)
	logger.Debug("CreateContractTx repository started")

	entity, err := scanContractRow(
		tx.QueryRow(
			ctx,
			createContract,
			t.ContractNumber,
			t.CompanyID,
			string(t.Status),
			t.StartDate,
			t.EndDate,
		))
	if err != nil {
		logger.Error(errSelectEntityFailed, slog.Any("error", err))
		return &domain.ContractEntity{}, fmt.Errorf("CreateContractTx: %w", err)
	}

	logger.Debug("CreateContractTx completed", slog.String("contract_id", entity.ID.String()))
	return entity, nil
}
