package usecase

import (
	"context"
	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

type BaseContractUseCase interface {
	GetContractByID(ctx context.Context, tx pgx.Tx, repo storage.BaseTxContractRepository, request *model.GetContractByIDRequest) (*model.Contract, error)
}

type ContractUseCase struct {
}

func NewContractUseCase() *ContractUseCase {
	return &ContractUseCase{}
}
