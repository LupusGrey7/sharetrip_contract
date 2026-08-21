package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
)

func (s *OfferingService) UpsertContractServices(
	ctx context.Context,
	input *domain.UpsertContractServicesInput,
) (*domain.UpsertContractServicesOutput, error) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("service", "OfferingService"),
		slog.String("operation", "UpsertContractServices"),
		slog.Int("contract_id", input.ContractID),
	)
	logger.Debug("upsert contract services started")

	res, err := tx(ctx, s.pool, func(pgTx pgx.Tx) (*domain.UpsertContractServicesOutput, error) {
		return s.useCase.UpsertContractServices(
			ctx,
			pgTx,
			s.contractRepo,
			s.offeringRepo,
			s.linkRepo,
			input,
		)
	})
	if err != nil {
		return nil, fmt.Errorf("UpsertContractServices: %w", err)
	}

	logger.Debug("upsert contract services completed")
	return res, nil
}
