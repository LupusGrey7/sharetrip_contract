package storage

import (
	"context"
	"fmt"
	"job4j/sharetrip-contract/internal/contract/model"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	getContractByID = `
	SELECT id, 
		   number,
		   date,
		   status,
		   total_amount 
	FROM contracts 
	WHERE id = $1;
	`
	createContract = `
	INSERT INTO contracts (number, date, status, total_amount) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, number, date, status, total_amount;
	`
	updateContract = `
	UPDATE contracts 
	SET number = $1,
		date = $2,
		status = $3,
		total_amount = $4 
	WHERE id = $5
	RETURNING id, number, date, status, total_amount;
	`
	changeContractStatus = `
	UPDATE contracts 
	SET status = $1
	WHERE id = $2
	RETURNING id, number, date, status, total_amount;
	`
)

type BaseTxContractRepository interface {
	GetContractByIDTx(ctx context.Context, tx pgx.Tx, id int) (*model.Contract, error)
	CreateContractTx(ctx context.Context, tx pgx.Tx, contract *model.Contract) (*model.Contract, error)
	UpdateContractTx(ctx context.Context, tx pgx.Tx, contract *model.Contract) (*model.Contract, error)
	ChangeContractStatusTx(ctx context.Context, tx pgx.Tx, id int, status model.ContractStatus) (*model.Contract, error)
}

type ContractRepository struct {
	pool *pgxpool.Pool
}

func NewContractRepository(pool *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{pool: pool}
}

func (r *ContractRepository) GetContractByIDTx(
	ctx context.Context,
	tx pgx.Tx,
	id int,
) (*model.Contract, error) {
	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "repository"),
		slog.String("repository", "ContractRepository.GetContractByIDTx"),
		slog.Int("contract_id", id),
	)
	logger.Debug("GetContractByIDTx repository started")

	var contract model.Contract
	err := tx.QueryRow(ctx, getContractByID, id).Scan(&contract.ID, &contract.Number, &contract.Date, &contract.Status, &contract.TotalAmount)
	if err != nil {
		logger.Error("GetContractByIDTx repository failed", slog.Any("error", err))
		return nil, err
	}
	if contract.ID == 0 {
		return nil, ErrContractNotFound
	}

	slog.Debug("GetContractByIDTx success")
	return &contract, nil
}

func (r *ContractRepository) CreateContractTx(
	ctx context.Context,
	tx pgx.Tx,
	contract *model.Contract,
) (*model.Contract, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "repository"),
		slog.String("repository", "CreateContractTx"),
		slog.Int("contract_id", contract.ID),
	)
	logger.Debug("CreateContractTx repository started")

	var newContract model.Contract
	var values []interface{}
	values = append(values, contract.Number, contract.Date, contract.Status, contract.TotalAmount)
	var query = createContract

	err := tx.QueryRow(ctx, query, values...).Scan(
		&newContract.ID,
		&newContract.Number,
		&newContract.Date,
		&newContract.Status,
		&newContract.TotalAmount,
	)

	if err != nil {
		logger.Error("CreateContractTx repository failed", slog.Any("error", err))
		return nil, err
	}

	logger.Debug("CreateContractTx repository completed")
	return contract, nil
}

func (r *ContractRepository) UpdateContractTx(
	ctx context.Context,
	tx pgx.Tx,
	contract *model.Contract,
) (*model.Contract, error) {
	return nil, fmt.Errorf("UpdateContractTx: not implemented")
}

func (r *ContractRepository) ChangeContractStatusTx(
	ctx context.Context,
	tx pgx.Tx,
	id int,
	status model.ContractStatus,
) (*model.Contract, error) {
	return nil, fmt.Errorf("ChangeContractStatusTx: not implemented")
}
