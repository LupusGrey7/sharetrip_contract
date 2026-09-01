package usecase

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/domain"
)

func (u *HealthcheckUseCase) GetHealthcheckInfo(ctx context.Context) (*domain.HealthcheckOutput, error) {
	if err := u.repo.CheckHealth(ctx); err != nil {
		return nil, fmt.Errorf("healthcheck: %w", err)
	}
	return &domain.HealthcheckOutput{
		Status:  "ok",
		Message: "application and database are healthy",
	}, nil
}
