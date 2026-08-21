package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

func (u *OfferingUseCase) UpsertContractServices(
	ctx context.Context,
	tx pgx.Tx,
	contractRepo storage.BaseTxContractRepository,
	offeringRepo storage.BaseTxOfferingRepository,
	linkRepo storage.BaseTxContractOfferingRepository,
	input *domain.UpsertContractServicesInput,
) (*domain.UpsertContractServicesOutput, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "useCase"),
		slog.String("useCase", "UpsertContractServices"),
		slog.Int("contract_id", input.ContractID),
	)
	logger.Debug("UpsertContractServices started")

	// 1) lock parent contract
	_, err := contractRepo.GetContractByIDForUpdateTx(ctx, tx, input.ContractID)
	if err != nil {
		if errors.Is(err, storage.ErrContractNotFound) {
			return nil, ErrContractNotFound
		}
		return nil, err
	}

	// 2) dictionary codes must exist (do not upsert dictionary)
	codes := make([]string, 0, len(input.Services))
	for _, item := range input.Services {
		codes = append(codes, item.ServiceCode)
	}
	missing, err := offeringRepo.ExistServiceCodesTx(ctx, tx, codes)
	if err != nil {
		return nil, err
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrServiceNotFound, strings.Join(missing, ","))
	}

	// 3) upsert links in contract_services
	for _, item := range input.Services {
		if err := linkRepo.UpsertContractServiceTx(
			ctx,
			tx,
			input.ContractID,
			domain.ServiceItemInputToEntity(item),
		); err != nil {
			return nil, err
		}
	}

	// 4) return current list
	items, err := linkRepo.ListContractServicesTx(ctx, tx, input.ContractID)
	if err != nil {
		return nil, err
	}

	logger.Debug("UpsertContractServices completed", slog.Int("count", len(items)))
	return &domain.UpsertContractServicesOutput{
		ContractID: input.ContractID,
		Services:   domain.ServiceItemEntitiesToOutput(items),
	}, nil
}
