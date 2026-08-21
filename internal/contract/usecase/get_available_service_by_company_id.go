package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

// GetAvailableOfferingByCompanyID — three explicit checks → three different outcomes for the client.
//
//  1. check: company exists = at least one contract with this company_id     → and ErrCompanyNotFound
//  2. check: service exists in the dictionary             → and ErrServiceNotFound (offering)
//  3. final check: active + contract_services        → 200 + allowed/reason (not 404)
func (u *CompanyUseCase) GetAvailableOfferingByCompanyID(
	ctx context.Context,
	tx pgx.Tx,
	offeringRepo storage.BaseTxOfferingRepository,
	companyRepo storage.BaseCompanyRepository,
	input *domain.GetAvailableOfferingByCompanyIDInput,
) (*domain.AvailabilityOutput, error) {
	log := slog.With(
		slog.String("company_id", strconv.Itoa(input.CompanyID)),
		slog.String("service_code", string(input.ServiceCode)),
	)
	log.Debug("get available offering by company id started")

	// 1)check: company exists = at least one contract with this company_id
	if err := companyRepo.IsCompanyKnownByIDTx(ctx, tx, input.CompanyID); err != nil {
		if errors.Is(err, storage.ErrCompanyNotFound) {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("IsCompanyKnownByIDTx: %w", err)
	}

	// 2)check: service exists in the dictionary
	if err := offeringRepo.IsOfferingExistsByCodeTx(ctx, tx, string(input.ServiceCode)); err != nil {
		if errors.Is(err, storage.ErrOfferingNotFound) {
			return nil, ErrServiceNotFound
		}
		return nil, fmt.Errorf("IsOfferingExistsByCodeTx: %w", err)
	}

	// 3)final check: active + contract_services
	entity, err := companyRepo.GetAvailableOfferingByCompanyIDTx(ctx, tx, input.CompanyID, input.ServiceCode)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableOfferingByCompanyIDTx: %w", err)
	}

	log.Debug("get available offering by company id completed", slog.Bool("allowed", entity.Allowed))
	return domain.AvailabilityEntityToOutput(entity), nil
}
