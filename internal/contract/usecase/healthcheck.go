package usecase

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/storage"
)

type BaseHealthcheckUseCase interface {
	GetHealthcheckInfo(ctx context.Context) (*domain.HealthcheckOutput, error)
}

type HealthcheckUseCase struct {
	repo storage.HealthRepository
}

func NewHealthcheckUseCase(repo storage.HealthRepository) *HealthcheckUseCase {
	return &HealthcheckUseCase{repo: repo}
}
