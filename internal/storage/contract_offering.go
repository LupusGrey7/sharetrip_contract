package storage

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseTxContractOfferingRepository — связи contract_services (не company_services).
type BaseTxContractOfferingRepository interface {
	UpsertContractServiceTx(ctx context.Context, tx pgx.Tx, contractID int, item model.ServiceItem) error
	ListContractServicesTx(ctx context.Context, tx pgx.Tx, contractID int) ([]model.ServiceItem, error)
}

type ContractOfferingRepository struct {
	pool *pgxpool.Pool
}

func NewContractOfferingRepository(pool *pgxpool.Pool) *ContractOfferingRepository {
	return &ContractOfferingRepository{pool: pool}
}
