package storage

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
)

const (
	getContractByID = `
SELECT id, contract_number, company_id, status_id, start_date, end_date, created_at, updated_at
FROM contract_management.contracts
WHERE id = $1
`

	getContractByIDForUpdate = getContractByID + ` FOR UPDATE`
)

func (r *ContractRepository) GetContractByIDTx(
	ctx context.Context,
	tx pgx.Tx,
	id int,
) (*domain.ContractEntity, error) {
	return r.scanContract(ctx, tx, getContractByID, id)
}

func (r *ContractRepository) GetContractByIDForUpdateTx(
	ctx context.Context,
	tx pgx.Tx,
	id int,
) (*domain.ContractEntity, error) {
	return r.scanContract(ctx, tx, getContractByIDForUpdate, id)
}

func (r *ContractRepository) scanContract(
	ctx context.Context,
	tx pgx.Tx,
	query string,
	id int,
) (*domain.ContractEntity, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "repository"),
		slog.String("repository", "ContractRepository"),
		slog.Int("contract_id", id),
	)

	contract, err := scanContractRow(tx.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrContractNotFound
		}
		logger.Error("GetContract failed", slog.Any("error", err))
		return nil, err
	}
	return contract, nil
}

func scanContractRow(row pgx.Row) (*domain.ContractEntity, error) {
	var c domain.ContractEntity
	var status string
	err := row.Scan(
		&c.ID,
		&c.ContractNumber,
		&c.CompanyID,
		&status,
		&c.StartDate,
		&c.EndDate,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	c.Status = domain.ContractStatus(status)
	return &c, nil
}
