package usecase

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

type BaseContractUseCase interface {
	GetContractByID(ctx context.Context, tx pgx.Tx, repo storage.BaseTxContractRepository, input *domain.GetContractByIDInput) (*domain.ContractOutput, error)
	CreateContract(ctx context.Context, tx pgx.Tx, repo storage.BaseTxContractRepository, input *domain.CreateContractInput) (*domain.ContractOutput, error)
}

type ContractUseCase struct{}

func NewContractUseCase() *ContractUseCase {
	return &ContractUseCase{}
}
