package service

import (
	"context"
	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/usecase"
)

type BaseHealthcheck interface {
	GetHealthcheckInfo(ctx context.Context) (*model.GetHealthcheckInfoResponse, error)
}

type HealthcheckService struct {
	useCase usecase.BaseInfoUseCase
}

func NewHealthcheckService(useCase usecase.BaseInfoUseCase) *HealthcheckService {
	return &HealthcheckService{useCase: useCase}
}
