package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

// GetAvailableOfferingByCompanyID — три явные проверки → три разных исхода для клиента.
//
//  1. company известна (есть договор с company_id)     → иначе ErrCompanyNotFound
//  2. service_code есть в словаре services             → иначе ErrServiceNotFound (offering)
//  3. финальный join: active + contract_services        → 200 + allowed/reason (не 404)
func (u *CompanyUseCase) GetAvailableOfferingByCompanyID(
	ctx context.Context,
	tx pgx.Tx,
	offeringRepo storage.BaseTxOfferingRepository,
	companyRepo storage.BaseCompanyRepository,
	req *model.GetAvailableOfferingByCompanyIDRequest,
) (*model.GetAvailableResultResponse, error) {
	log := slog.With(
		slog.String("company_id", strconv.Itoa(req.CompanyID)),
		slog.String("service_code", string(req.ServiceCode)),
	)
	log.Debug("get available offering by company id started")

	if req == nil || req.CompanyID < 1 || req.ServiceCode == "" {
		return nil, fmt.Errorf("%w: company_id and service_code required", ErrInvalidRequest)
	}

	// 1) компания «есть у нас» = хотя бы один contract с этим company_id
	if err := companyRepo.IsCompanyKnownByIDTx(ctx, tx, req.CompanyID); err != nil {
		if errors.Is(err, storage.ErrCompanyNotFound) {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("IsCompanyKnownByIDTx: %w", err)
	}

	// 2) услуга есть в словаре
	if err := offeringRepo.IsOfferingExistsByCodeTx(ctx, tx, string(req.ServiceCode)); err != nil {
		if errors.Is(err, storage.ErrOfferingNotFound) {
			return nil, ErrServiceNotFound
		}
		return nil, fmt.Errorf("IsOfferingExistsByCodeTx: %w", err)
	}

	// 3) бизнес-ответ: можно ли пользоваться прямо сейчас
	resp, err := companyRepo.GetAvailableOfferingByCompanyIDTx(ctx, tx, req.CompanyID, req.ServiceCode)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableOfferingByCompanyIDTx: %w", err)
	}

	log.Debug("get available offering by company id completed", slog.Bool("allowed", resp.Allowed))
	return resp, nil
}
