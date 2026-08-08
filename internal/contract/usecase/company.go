package usecase

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/model"
)

type BaseCompanyUseCase interface {
	GetAvailableServiceByCompanyID(ctx context.Context, req *model.GetAvailableServiceByCompanyIDRequest) (*model.GetAvailableResultResponse, error)
}

type CompanyUseCase struct{}

func NewCompanyUseCase() *CompanyUseCase {
	return &CompanyUseCase{}
}
