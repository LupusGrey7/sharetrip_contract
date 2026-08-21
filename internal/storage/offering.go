package storage

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseTxOfferingRepository — словарь services.
type BaseTxOfferingRepository interface {
	ExistServiceCodesTx(ctx context.Context, tx pgx.Tx, codes []string) (missing []string, err error)
	GetOfferingByCodeTx(ctx context.Context, tx pgx.Tx, code string) (*domain.OfferingEntity, error)
	IsOfferingExistsByCodeTx(ctx context.Context, tx pgx.Tx, code string) error
}

type OfferingRepository struct {
	pool *pgxpool.Pool
}

func NewOfferingRepository(pool *pgxpool.Pool) *OfferingRepository {
	return &OfferingRepository{pool: pool}
}
