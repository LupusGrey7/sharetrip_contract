package usecase

import (
	"context"
	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

// BaseOfferingUseCase is the interface for the offering use case.
type BaseOfferingUseCase interface {
	UpsertContractServices(ctx context.Context, tx pgx.Tx, repo storage.BaseTxOfferingRepository, request *model.UpsertContractOfferingsRequest) (*model.UpsertContractOfferingsResponse, error)
	GetOfferingByCode(ctx context.Context, tx pgx.Tx, repo storage.BaseTxOfferingRepository, request *model.GetOfferingByCodeRequest) (*model.GetOfferingByCodeResponse, error)
}

type OfferingUseCase struct {
}

func NewOfferingUseCase() *OfferingUseCase {
	return &OfferingUseCase{}
}
