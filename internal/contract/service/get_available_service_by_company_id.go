package service

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/model"
)

func (s *CompanyService) GetAvailableServiceByCompanyID(
	ctx context.Context,
	req *model.GetAvailableServiceByCompanyIDRequest,
) (*model.GetAvailableResultResponse, error) {
	resp, err := s.useCase.GetAvailableServiceByCompanyID(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableServiceByCompanyID: %w", err)
	}
	return resp, nil
}
