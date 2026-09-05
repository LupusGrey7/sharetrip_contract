package usecase

import (
	"context"
	"errors"
	"log/slog"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/observability/logctx"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

func (u *ContractUseCase) GetActiveContractByCompanyID(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	input *domain.GetActiveContractByCompanyIDInput,
) (*domain.ContractOutput, error) {
	ctxSpc, span := otel.Tracer("ContractUseCase").Start(ctx, "ContractUseCase.GetActiveContractByCompanyID")
	defer span.End()

	logger := logctx.Logger(ctxSpc).With(
		slog.String("layer", "useCase"),
		slog.String("useCase", "GetActiveContractByCompanyID"),
		slog.Int("company_id", input.CompanyID),
	)
	logger.Debug("GetActiveContractByCompanyID started")

	entity, err := repo.GetActiveContractByCompanyIDTx(ctxSpc, tx, input.CompanyID)
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
