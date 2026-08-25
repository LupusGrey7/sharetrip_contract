package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

func (u *ContractUseCase) GetActiveContractByCompanyID(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	input *domain.GetActiveContractByCompanyIDInput,
) (*domain.ContractOutput, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "useCase"),
		slog.String("useCase", "GetActiveContractByCompanyID"),
		slog.Int("company_id", input.CompanyID),
	)
	logger.Debug("GetActiveContractByCompanyID started")

	entity, err := repo.GetActiveContractByCompanyIDTx(ctx, tx, input.CompanyID)
	if err != nil {
		logger.Error("GetActiveContractByCompanyID failed", slog.Any("error", err))
		if errors.Is(err, storage.ErrContractNotFound) {
			return nil, ErrContractNotFound
		}
		return nil, err
	}

	logger.Debug("GetActiveContractByCompanyID completed")
	return domain.ContractEntityToOutput(entity), nil
}
