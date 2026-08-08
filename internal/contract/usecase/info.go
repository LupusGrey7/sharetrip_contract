package usecase

import (
	"context"
	"fmt"

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

func (u *HealthcheckUseCase) GetHealthcheckInfo(ctx context.Context) (*model.GetHealthcheckInfoResponse, error) {
	if err := u.repo.CheckHealth(ctx); err != nil {
		return nil, fmt.Errorf("healthcheck: %w", err)
	}
	return &model.GetHealthcheckInfoResponse{
		Status:  "ok",
		Message: "application and database are healthy",
	}, nil
}
