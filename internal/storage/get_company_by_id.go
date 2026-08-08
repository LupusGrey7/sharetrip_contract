package storage

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5"
)

// «Компания известна системе» = есть хотя бы один договор с этим company_id.
// Таблицы companies в миграциях нет — company_id внешний ID.
const isCompanyKnownByIDQuery = `
SELECT EXISTS(
	SELECT 1
	FROM contract_management.contracts
	WHERE company_id = $1
)
`

func (r *CompanyRepository) IsCompanyKnownByIDTx(ctx context.Context, tx pgx.Tx, companyID int) error {
	log := slog.With(slog.String("company_id", strconv.Itoa(companyID)))
	log.Debug("is company known by id started")

	var exists bool
	if err := tx.QueryRow(ctx, isCompanyKnownByIDQuery, companyID).Scan(&exists); err != nil {
		return fmt.Errorf("IsCompanyKnownByIDTx: %w", err)
	}
	if !exists {
		return ErrCompanyNotFound
	}

	log.Debug("is company known by id success")
	return nil
}
