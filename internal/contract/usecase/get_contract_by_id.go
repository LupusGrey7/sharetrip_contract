// scenario: get contract by id - business logic + repository
package usecase

import (
	"context"
	"errors"
	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func (u *ContractUseCase) GetContractByID(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	req *model.GetContractByIDRequest,
) (*model.Contract, error) {
	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "useCase"),
		slog.String("useCase", "ContractUseCase.GetContractByID"),
		slog.Int("contract_id", req.ContractID),
	)
	logger.Debug("GetContractByID useCase started")

	contract, err := repo.GetContractByIDTx(ctx, tx, req.ContractID)
	if err != nil {
		logger.Error("GetContractByID useCase failed", slog.Any("error", err))
		if errors.Is(err, storage.ErrContractNotFound) {
			return nil, ErrContractNotFound
		}
		// If this is not a ErrEntityNotFound, This means it's a system error (500 error)
		return nil, err
	}

	logger.Debug("GetContractByID useCase completed")
	return contract, nil
}
