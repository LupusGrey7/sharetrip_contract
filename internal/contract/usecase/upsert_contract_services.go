package usecase

import (
	"context"
	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

func (u *OfferingUseCase) UpsertContractServices(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxOfferingRepository,
	request *model.UpsertContractOfferingsRequest,
) (*model.UpsertContractOfferingsResponse, error) {
	return nil, nil
}

func (u *OfferingUseCase) GetOfferingByCode(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxOfferingRepository,
	request *model.GetOfferingByCodeRequest,
) (*model.GetOfferingByCodeResponse, error) {
	return nil, nil
}
