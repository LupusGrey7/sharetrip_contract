package storage

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseTxContractRepository interface {
	GetContractByIDTx(ctx context.Context, tx pgx.Tx, id int) (*model.Contract, error)
	GetContractByIDForUpdateTx(ctx context.Context, tx pgx.Tx, id int) (*model.Contract, error)
	CreateContractTx(ctx context.Context, tx pgx.Tx, contract *model.Contract) (*model.Contract, error)
}

type ContractRepository struct {
	pool *pgxpool.Pool
}

func NewContractRepository(pool *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{pool: pool}
}
