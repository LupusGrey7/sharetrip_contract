package storage

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
)

const (
	existServiceCodes = `
SELECT service_code
FROM contract_management.services
WHERE service_code = ANY($1)`

	getOfferingByCode = `
SELECT service_code, description, is_active, created_at, updated_at
FROM contract_management.services
WHERE service_code = $1`
)

func (r *OfferingRepository) ExistServiceCodesTx(
	ctx context.Context,
	tx pgx.Tx,
	codes []string,
) ([]string, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	rows, err := tx.Query(ctx, existServiceCodes, codes)
	if err != nil {
		return nil, fmt.Errorf("ExistServiceCodesTx: %w", err)
	}
	defer rows.Close()

	found := make(map[string]struct{}, len(codes))
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, fmt.Errorf("ExistServiceCodesTx scan: %w", err)
		}
		found[code] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var missing []string
	for _, code := range codes {
		if _, ok := found[code]; !ok {
			missing = append(missing, code)
		}
	}
	return missing, nil
}

func (r *OfferingRepository) GetOfferingByCodeTx(
	ctx context.Context,
	tx pgx.Tx,
	code string,
) (*domain.OfferingEntity, error) {
	var o domain.OfferingEntity
	err := tx.QueryRow(ctx, getOfferingByCode, code).Scan(
		&o.ServiceCode,
		&o.Description,
		&o.IsActive,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrOfferingNotFound
		}
		return nil, fmt.Errorf("GetOfferingByCodeTx: %w", err)
	}
	return &o, nil
}
