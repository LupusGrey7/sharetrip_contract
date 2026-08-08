package service

import (
	"context"
	"fmt"
	"job4j/sharetrip-contract/internal/contract/model"
	"log/slog"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func (s *OfferingService) UpsertContractOfferings(
	ctx context.Context,
	request *model.UpsertContractOfferingsRequest,
) (*model.UpsertContractOfferingsResponse, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("service", "OfferingService"),
		slog.String("operation", "UpsertContractOfferings"),
		slog.String("contract_id", strconv.Itoa(request.ContractID)),
	)
	logger.Debug("upsert contract offerings started")

	resp, err := tx(ctx, s.pool, func(pgTx pgx.Tx) (*model.UpsertContractOfferingsResponse, error) {
		return s.useCase.UpsertContractServices(ctx, pgTx, s.repo, request)
	})
	if err != nil {
		logger.Error("upsert contract offerings failed", slog.Any("error", err))
		return nil, fmt.Errorf("UpsertContractOfferings: %w", err)
	}

	logger.Debug("upsert contract offerings completed")
	return resp, nil
}
