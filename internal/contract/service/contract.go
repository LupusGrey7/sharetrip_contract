package service

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Contract — контракт сервисного слоя (методы в отдельных файлах).
type Contract interface {
	GetContractByID(ctx context.Context, input *domain.GetContractByIDInput) (*domain.ContractOutput, error)
	GetActiveContractByCompanyID(ctx context.Context, input *domain.GetActiveContractByCompanyIDInput) (*domain.ContractOutput, error)
	CreateContract(ctx context.Context, input *domain.CreateContractInput) (*domain.ContractOutput, error)
}

type ContractService struct {
	pool    *pgxpool.Pool
	repo    storage.BaseTxContractRepository
	useCase usecase.BaseContractUseCase
}

func NewContractService(
	pool *pgxpool.Pool,
	repo storage.BaseTxContractRepository,
	contractUseCase usecase.BaseContractUseCase,
) *ContractService {
	return &ContractService{
		pool:    pool,
		repo:    repo,
		useCase: contractUseCase,
	}
}
