package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseTxOfferingRepository — словарь services (для availability).
type BaseTxOfferingRepository interface {
	IsOfferingExistsByCodeTx(ctx context.Context, tx pgx.Tx, code string) error
}

type OfferingRepository struct {
	pool *pgxpool.Pool
}

func NewOfferingRepository(pool *pgxpool.Pool) *OfferingRepository {
	return &OfferingRepository{pool: pool}
}
