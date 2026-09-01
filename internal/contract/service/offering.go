package service

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Offering — контракт сервисного слоя (методы в отдельных файлах).
type Offering interface {
	UpsertContractServices(ctx context.Context, input *domain.UpsertContractServicesInput) (*domain.UpsertContractServicesOutput, error)
}

type OfferingService struct {
	pool         *pgxpool.Pool
	contractRepo storage.BaseTxContractRepository
	offeringRepo storage.BaseTxOfferingRepository
	linkRepo     storage.BaseTxContractOfferingRepository
	useCase      usecase.BaseOfferingUseCase
}

func NewOfferingService(
	pool *pgxpool.Pool,
	contractRepo storage.BaseTxContractRepository,
	offeringRepo storage.BaseTxOfferingRepository,
	linkRepo storage.BaseTxContractOfferingRepository,
	useCase usecase.BaseOfferingUseCase,
) *OfferingService {
	return &OfferingService{
		pool:         pool,
		contractRepo: contractRepo,
		offeringRepo: offeringRepo,
		linkRepo:     linkRepo,
		useCase:      useCase,
	}
}
