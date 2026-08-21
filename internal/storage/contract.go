package storage

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseTxContractRepository interface {
	GetContractByIDTx(ctx context.Context, tx pgx.Tx, id int) (*domain.ContractEntity, error)
	GetContractByIDForUpdateTx(ctx context.Context, tx pgx.Tx, id int) (*domain.ContractEntity, error)
	CreateContractTx(ctx context.Context, tx pgx.Tx, entity *domain.ContractEntity) (*domain.ContractEntity, error)
}

type ContractRepository struct {
	pool *pgxpool.Pool
}

func NewContractRepository(pool *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{pool: pool}
}
