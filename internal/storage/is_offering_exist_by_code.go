package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

const isOfferingExistsByCodeQuery = `
SELECT EXISTS(
	SELECT 1
	FROM contract_management.services
	WHERE service_code = $1
)
`

func (r *OfferingRepository) IsOfferingExistsByCodeTx(ctx context.Context, tx pgx.Tx, code string) error {
	log := slog.With(slog.String("offering_code", code))
	log.Debug("is offering exists by code started")

	var exists bool
	if err := tx.QueryRow(ctx, isOfferingExistsByCodeQuery, code).Scan(&exists); err != nil {
		return fmt.Errorf("IsOfferingExistsByCodeTx: %w", err)
	}
	if !exists {
		return ErrOfferingNotFound
	}

	log.Debug("is offering exists by code success")
	return nil
}
