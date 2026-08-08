package storage

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/jackc/pgx/v5"
)

const (
	upsertContractService = `
INSERT INTO contract_management.contract_services (contract_id, service_code, is_enabled)
VALUES ($1, $2, $3)
ON CONFLICT (contract_id, service_code)
DO UPDATE SET
	is_enabled = EXCLUDED.is_enabled,
	updated_at = NOW()`

	listContractServices = `
SELECT service_code, is_enabled
FROM contract_management.contract_services
WHERE contract_id = $1
ORDER BY service_code`
)

func (r *ContractOfferingRepository) UpsertContractServiceTx(
	ctx context.Context,
	tx pgx.Tx,
	contractID int,
	item model.ServiceItem,
) error {
	_, err := tx.Exec(ctx, upsertContractService, contractID, item.ServiceCode, item.IsEnabled)
	if err != nil {
		return fmt.Errorf("UpsertContractServiceTx: %w", err)
	}
	return nil
}

func (r *ContractOfferingRepository) ListContractServicesTx(
	ctx context.Context,
	tx pgx.Tx,
	contractID int,
) ([]model.ServiceItem, error) {
	rows, err := tx.Query(ctx, listContractServices, contractID)
	if err != nil {
		return nil, fmt.Errorf("ListContractServicesTx: %w", err)
	}
	defer rows.Close()

	var items []model.ServiceItem
	for rows.Next() {
		var item model.ServiceItem
		if err := rows.Scan(&item.ServiceCode, &item.IsEnabled); err != nil {
			return nil, fmt.Errorf("ListContractServicesTx scan: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if items == nil {
		items = []model.ServiceItem{}
	}
	return items, nil
}
