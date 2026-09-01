package storage

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseCompanyRepository interface {
	IsCompanyKnownByIDTx(ctx context.Context, tx pgx.Tx, companyID int) error
	GetAvailableOfferingByCompanyIDTx(ctx context.Context, tx pgx.Tx, companyID int, serviceCode domain.ServiceCode) (*domain.AvailabilityEntity, error)
}

type CompanyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{pool: pool}
}
