package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

func (u *ContractUseCase) GetContractByID(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	req *model.GetContractByIDRequest,
) (*model.Contract, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "useCase"),
		slog.String("useCase", "GetContractByID"),
		slog.Int("contract_id", req.ContractID),
	)
	logger.Debug("GetContractByID started")

	contract, err := repo.GetContractByIDTx(ctx, tx, req.ContractID)
	if err != nil {
		logger.Error("GetContractByID failed", slog.Any("error", err))
		if errors.Is(err, storage.ErrContractNotFound) {
			return nil, ErrContractNotFound
		}
		return nil, err
	}

	logger.Debug("GetContractByID completed")
	return contract, nil
}
