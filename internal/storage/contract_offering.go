package storage

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseTxContractOfferingRepository — связи contract_services (не company_services).
type BaseTxContractOfferingRepository interface {
	UpsertContractServiceTx(ctx context.Context, tx pgx.Tx, contractID int, item domain.ServiceItemEntity) error
	ListContractServicesTx(ctx context.Context, tx pgx.Tx, contractID int) ([]domain.ServiceItemEntity, error)
}

type ContractOfferingRepository struct {
	pool *pgxpool.Pool
}

func NewContractOfferingRepository(pool *pgxpool.Pool) *ContractOfferingRepository {
	return &ContractOfferingRepository{pool: pool}
}
