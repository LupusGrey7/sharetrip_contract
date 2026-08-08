package storage

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5"
)

// "Company is known to the system" = there is at least one contract with this company_id.
// The companies table is not in the migrations - company_id is an external ID.
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
