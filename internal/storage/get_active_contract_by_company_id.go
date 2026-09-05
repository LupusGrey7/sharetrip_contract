package storage

import (
	"context"
	"errors"
	"log/slog"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/observability/logctx"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

const getActiveContractByCompanyID = `
SELECT id, contract_number, company_id, status_id, start_date, end_date, created_at, updated_at
FROM contract_management.contracts
WHERE company_id = $1 AND status_id = 'active'
LIMIT 1
`

func (r *ContractRepository) GetActiveContractByCompanyIDTx(
	ctx context.Context,
	tx pgx.Tx,
	companyID int,
) (*domain.ContractEntity, error) {
	ctxSpc, span := otel.Tracer("ContractRepository").Start(ctx, "ContractRepository.GetActiveContractByCompanyIDTx")
	defer span.End()

	logger := logctx.Logger(ctxSpc).With(
		slog.String("layer", "repository"),
		slog.String("repository", "GetActiveContractByCompanyIDTx"),
		slog.Int("company_id", companyID),
	)

	contract, err := scanContractRow(tx.QueryRow(ctxSpc, getActiveContractByCompanyID, companyID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrContractNotFound
		}
		logger.Error("GetActiveContractByCompanyIDTx failed", slog.Any("error", err))
		return nil, err
	}
	return contract, nil
}
