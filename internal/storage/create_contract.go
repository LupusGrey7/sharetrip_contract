package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/jackc/pgx/v5"
)

const createContract = `
INSERT INTO contract_management.contracts (
	contract_number, company_id, status_id, start_date, end_date
) VALUES ($1, $2, $3, $4, $5)
RETURNING id, contract_number, company_id, status_id, start_date, end_date, created_at, updated_at`

func (r *ContractRepository) CreateContractTx(
	ctx context.Context,
	tx pgx.Tx,
	contract *model.Contract,
) (*model.Contract, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "repository"),
		slog.String("repository", "CreateContractTx"),
		slog.Int("company_id", contract.CompanyID),
	)
	logger.Debug("CreateContractTx started")

	created, err := scanContractRow(tx.QueryRow(
		ctx,
		createContract,
		contract.ContractNumber,
		contract.CompanyID,
		string(contract.Status),
		contract.StartDate,
		contract.EndDate,
	))
	if err != nil {
		logger.Error("CreateContractTx failed", slog.Any("error", err))
		return nil, fmt.Errorf("CreateContractTx: %w", err)
	}

	logger.Debug("CreateContractTx completed", slog.Int("contract_id", created.ID))
	return created, nil
}
