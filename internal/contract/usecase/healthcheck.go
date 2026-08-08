package usecase

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"
)

type BaseHealthcheckUseCase interface {
	GetHealthcheckInfo(ctx context.Context) (*model.GetHealthcheckInfoResponse, error)
}

type HealthcheckUseCase struct {
	repo storage.HealthRepository
}

func NewHealthcheckUseCase(repo storage.HealthRepository) *HealthcheckUseCase {
	return &HealthcheckUseCase{repo: repo}
}
