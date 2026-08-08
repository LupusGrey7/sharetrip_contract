package usecase

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/model"
)

func (u *HealthcheckUseCase) GetHealthcheckInfo(ctx context.Context) (*model.GetHealthcheckInfoResponse, error) {
	if err := u.repo.CheckHealth(ctx); err != nil {
		return nil, fmt.Errorf("healthcheck: %w", err)
	}
	return &model.GetHealthcheckInfoResponse{
		Status:  "ok",
		Message: "application and database are healthy",
	}, nil
}
