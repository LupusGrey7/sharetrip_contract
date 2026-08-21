package usecase

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

type BaseOfferingUseCase interface {
	UpsertContractServices(
		ctx context.Context,
		tx pgx.Tx,
		contractRepo storage.BaseTxContractRepository,
		offeringRepo storage.BaseTxOfferingRepository,
		linkRepo storage.BaseTxContractOfferingRepository,
		input *domain.UpsertContractServicesInput,
	) (*domain.UpsertContractServicesOutput, error)
}

type OfferingUseCase struct{}

func NewOfferingUseCase() *OfferingUseCase {
	return &OfferingUseCase{}
}
