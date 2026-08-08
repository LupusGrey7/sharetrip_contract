package service

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/jackc/pgx/v5"
)

func (s *CompanyService) GetAvailableOfferingByCompanyID(
	ctx context.Context,
	req *model.GetAvailableOfferingByCompanyIDRequest,
) (*model.GetAvailableResultResponse, error) {
	// pool == nil - only unit tests without DB (direct use case call).
	if s.pool == nil {
		res, err := s.useCase.GetAvailableOfferingByCompanyID(
			ctx, nil, s.offeringRepo, s.companyRepo, req,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableOfferingByCompanyID: %w", err)
		}
		return res, nil
	}

	res, err := tx(ctx, s.pool, func(pgTx pgx.Tx) (*model.GetAvailableResultResponse, error) {
		return s.useCase.GetAvailableOfferingByCompanyID(
			ctx,
			pgTx,
			s.offeringRepo,
			s.companyRepo,
			req,
		)
	})
	if err != nil {
		return nil, fmt.Errorf("GetAvailableOfferingByCompanyID: %w", err)
	}
	return res, nil
}
