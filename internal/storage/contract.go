package storage

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseTxContractRepository interface {
	GetActiveContractByCompanyIDTx(ctx context.Context, tx pgx.Tx, companyID int) (*domain.ContractEntity, error)
	CreateContractTx(ctx context.Context, tx pgx.Tx, entity *domain.ContractEntity) (*domain.ContractEntity, error)
}

type ContractRepository struct {
	pool *pgxpool.Pool
}

func NewContractRepository(pool *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{pool: pool}
}
