package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/observability/logctx"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

func (u *ContractUseCase) CreateContract(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	input *domain.CreateContractInput,
) (*domain.ContractOutput, error) {
	//tracing Jaeger
	ctxSpc, span := otel.Tracer("TripContractUseCase").Start(ctx, "TripContractUseCase.CreateContract")
	defer span.End()

	logger := logctx.Logger(ctxSpc).With(
		slog.String("layer", "useCase"),
		slog.String("useCase", "CreateContract"),
		slog.Int("company_id", input.CompanyID),
	)
	logger.Debug("create contract usecase started")

	if input.EndDate.Before(input.StartDate) {
		return nil, fmt.Errorf("%w: end_date before start_date", ErrInvalidRequest)
	}

	status := input.Status
	if status == "" {
		status = domain.ContractStatusDraft
	}

	number := input.ContractNumber
	if number == "" {
		number = fmt.Sprintf("C-%d-%d", input.CompanyID, time.Now().Unix())
	}

	entity, err := repo.CreateContractTx(ctx, tx, &domain.ContractEntity{
		ContractNumber: number,
		CompanyID:      input.CompanyID,
		Status:         status,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
	})
	if err != nil {
		logger.Error("repository create contract failed", slog.Any("error", err))
		return nil, err
	}

	logger.Debug("create contract usecase completed", slog.Int("contract_id", entity.ID.String())
	return domain.ContractEntityToOutput(entity), nil
}
