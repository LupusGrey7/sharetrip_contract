package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"job4j/sharetrip-contract/internal/contract/model"
)

func (s *HealthcheckService) GetHealthcheckInfo(ctx context.Context) (*model.GetHealthcheckInfoResponse, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("service", "HealthcheckService"),
		slog.String("operation", "GetHealthcheckInfo"),
	)
	logger.Debug("get healthcheck info started")

	resp, err := s.useCase.GetHealthcheckInfo(ctx)
	if err != nil {
		logger.Error("get healthcheck info failed", slog.Any("error", err))
		return nil, fmt.Errorf("GetHealthcheckInfo: %w", err)
	}

	logger.Debug("get healthcheck info completed")
	return resp, nil
}
