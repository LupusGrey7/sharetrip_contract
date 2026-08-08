package service

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/model"

	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func (s *OfferingService) UpsertContractServices(
	ctx context.Context,
	request *model.UpsertContractServicesRequest,
) (*model.UpsertContractServicesResponse, error) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("service", "OfferingService"),
		slog.String("operation", "UpsertContractServices"),
		slog.Int("contract_id", request.ContractID),
	)
	logger.Debug("upsert contract services started")

	res, err := tx(ctx, s.pool, func(pgTx pgx.Tx) (*model.UpsertContractServicesResponse, error) {
		return s.useCase.UpsertContractServices(
			ctx,
			pgTx,
			s.contractRepo,
			s.offeringRepo,
			s.linkRepo,
			request,
		)
	})
	if err != nil {
		return nil, fmt.Errorf("UpsertContractServices: %w", err)
	}

	logger.Debug("upsert contract services completed")
	return res, nil
}
