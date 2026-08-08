package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

func (u *OfferingUseCase) UpsertContractServices(
	ctx context.Context,
	tx pgx.Tx,
	contractRepo storage.BaseTxContractRepository,
	offeringRepo storage.BaseTxOfferingRepository,
	linkRepo storage.BaseTxContractOfferingRepository,
	req *model.UpsertContractServicesRequest,
) (*model.UpsertContractServicesResponse, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "useCase"),
		slog.String("useCase", "UpsertContractServices"),
		slog.Int("contract_id", req.ContractID),
	)
	logger.Debug("UpsertContractServices started")

	if req == nil || len(req.Services) == 0 {
		return nil, fmt.Errorf("%w: services required", ErrInvalidRequest)
	}

	// 1) lock parent contract
	_, err := contractRepo.GetContractByIDForUpdateTx(ctx, tx, req.ContractID)
	if err != nil {
		if errors.Is(err, storage.ErrContractNotFound) {
			return nil, ErrContractNotFound
		}
		return nil, err
	}

	// 2) dictionary codes must exist (do not upsert dictionary)
	codes := make([]string, 0, len(req.Services))
	for _, item := range req.Services {
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
	for _, item := range req.Services {
		if err := linkRepo.UpsertContractServiceTx(ctx, tx, req.ContractID, item); err != nil {
			return nil, err
		}
	}

	// 4) return current list
	items, err := linkRepo.ListContractServicesTx(ctx, tx, req.ContractID)
	if err != nil {
		return nil, err
	}

	logger.Debug("UpsertContractServices completed", slog.Int("count", len(items)))
	return &model.UpsertContractServicesResponse{
		ContractID: req.ContractID,
		Services:   items,
	}, nil
}
