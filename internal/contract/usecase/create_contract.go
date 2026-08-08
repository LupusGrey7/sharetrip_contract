package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

func (u *ContractUseCase) CreateContract(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	req *model.CreateContractRequest,
) (*model.Contract, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "useCase"),
		slog.String("useCase", "CreateContract"),
		slog.Int("company_id", req.CompanyID),
	)
	logger.Debug("CreateContract started")

	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("%w: end_date before start_date", ErrInvalidRequest)
	}

	status := req.Status
	if status == "" {
		status = model.ContractStatusDraft
	}

	number := req.ContractNumber
	if number == "" {
		number = fmt.Sprintf("C-%d-%d", req.CompanyID, time.Now().Unix())
	}

	resp, err := repo.CreateContractTx(ctx, tx, &model.Contract{
		ContractNumber: number,
		CompanyID:      req.CompanyID,
		Status:         status,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
	})
	if err != nil {
		logger.Error("CreateContract failed", slog.Any("error", err))
		return nil, err
	}

	logger.Debug("CreateContract completed", slog.Int("contract_id", resp.ID))
	return resp, nil
}
