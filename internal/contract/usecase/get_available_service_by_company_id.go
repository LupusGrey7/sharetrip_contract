package usecase

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/model"
)

// GetAvailableServiceByCompanyID — stub: допиши бизнес-правила (active contract + enabled service).
func (u *CompanyUseCase) GetAvailableServiceByCompanyID(
	ctx context.Context,
	req *model.GetAvailableServiceByCompanyIDRequest,
) (*model.GetAvailableResultResponse, error) {
	if req == nil || req.CompanyID < 1 || req.ServiceCode == "" {
		return nil, fmt.Errorf("%w: company_id and service_code required", ErrInvalidRequest)
	}

	// TODO: storage — active contract by company_id + contract_services.is_enabled
	return &model.GetAvailableResultResponse{
		CompanyID:   req.CompanyID,
		ServiceCode: req.ServiceCode,
		Allowed:     false,
		Reason:      "not implemented",
	}, nil
}
