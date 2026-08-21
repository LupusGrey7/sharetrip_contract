package service

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/usecase"
)

type BaseHealthcheck interface {
	GetHealthcheckInfo(ctx context.Context) (*domain.HealthcheckOutput, error)
}

type HealthcheckService struct {
	useCase usecase.BaseHealthcheckUseCase
}

func NewHealthcheckService(useCase usecase.BaseHealthcheckUseCase) *HealthcheckService {
	return &HealthcheckService{useCase: useCase}
}
