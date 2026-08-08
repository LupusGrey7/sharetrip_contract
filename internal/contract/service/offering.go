package service

import (
	"context"
	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OfferingService is responsible for managing commercial offerings (services) for a contract
type Offering interface {
	UpsertContractOfferings(ctx context.Context, request *model.UpsertContractOfferingsRequest) (*model.UpsertContractOfferingsResponse, error)
	GetOfferingByCode(ctx context.Context, request *model.GetOfferingByCodeRequest) (*model.GetOfferingByCodeResponse, error)
}

type OfferingService struct {
	pool    *pgxpool.Pool
	repo    storage.BaseTxOfferingRepository
	useCase usecase.BaseOfferingUseCase
}

func NewOfferingService(
	pool *pgxpool.Pool,
	repo storage.BaseTxOfferingRepository,
	useCase usecase.BaseOfferingUseCase,
) *OfferingService {
	return &OfferingService{
		pool:    pool,
		repo:    repo,
		useCase: useCase,
	}
}
