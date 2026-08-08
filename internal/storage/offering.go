package storage

import (
	"context"
	"fmt"
	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	getOfferingByID = `
	SELECT service_code,
		   description,
		   is_active,
		   created_at,
		   updated_at 
	FROM services 
	WHERE service_code = $1;
	`
	insertOffering = `
	INSERT INTO services (service_code, description, is_active) 
	VALUES ($1, $2, $3) 
	ON CONFLICT (service_code) DO NOTHING
	RETURNING service_code, description, is_active, created_at, updated_at;
	`

	updateOffering = `
	UPDATE services 
	SET description = $1,
		is_active = $2,
		updated_at = NOW()
	WHERE service_code = $3
	ON CONFLICT (service_code) DO NOTHING
	RETURNING service_code, description, is_active, created_at, updated_at;
	`

	getOfferingByCode = `
	SELECT service_code,
		   description,
		   is_active,
		   created_at,
		   updated_at 
	FROM services 
	WHERE service_code = $1;
	`

	getOfferingByIDForUpdate = `
	SELECT service_code,
		   description,
		   is_active,
		   created_at,
		   updated_at 
	FROM services 
	WHERE service_code = $1 FOR UPDATE;
	`
)

type BaseTxOfferingRepository interface {
	GetOfferingByIDTx(ctx context.Context, tx pgx.Tx, id int) (*model.Offering, error)
	UpsertOfferingTx(ctx context.Context, tx pgx.Tx, offering *model.Offering) error
	GetOfferingByCodeTx(ctx context.Context, tx pgx.Tx, code string) (*model.Offering, error)
}

type OfferingRepository struct {
	pool *pgxpool.Pool
}

func NewOfferingRepository(pool *pgxpool.Pool) *OfferingRepository {
	return &OfferingRepository{pool: pool}
}

func (r *OfferingRepository) GetOfferingByIDTx(
	ctx context.Context,
	tx pgx.Tx,
	id int,
) (*model.Offering, error) {
	return nil, fmt.Errorf("GetOfferingByIDTx: not implemented")
}

func (r *OfferingRepository) UpsertOfferingTx(
	ctx context.Context,
	tx pgx.Tx,
	offering *model.Offering,
) error {
	return fmt.Errorf("UpsertOfferingTx: not implemented")
}

func (r *OfferingRepository) GetOfferingByCodeTx(
	ctx context.Context,
	tx pgx.Tx,
	code string,
) (*model.Offering, error) {
	return nil, fmt.Errorf("GetOfferingByCodeTx: not implemented")
}
