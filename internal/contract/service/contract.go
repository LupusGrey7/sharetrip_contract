package service

import (
	"context"
	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ContractService is responsible for managing contracts
type Contract interface {
	GetContractByID(ctx context.Context, request *model.GetContractByIDRequest) (*model.GetContractByIDResponse, error)
	CreateContract(ctx context.Context, request *model.CreateContractRequest) error
	ChangeContractStatus(ctx context.Context, request *model.ChangeContractStatusRequest) (model.ChangeContractStatusResponse, error)
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
