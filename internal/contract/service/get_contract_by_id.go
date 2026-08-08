package service

import (
	"context"
	"fmt"
	"job4j/sharetrip-contract/internal/contract/model"
	"log/slog"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func (s *ContractService) GetContractByID(
	ctx context.Context,
	request *model.GetContractByIDRequest,
) (res *model.GetContractByIDResponse, err error) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("service", "ContractService"),
		slog.String("operation", "GetContractByID"),
		slog.String("contract_id", strconv.Itoa(request.ContractID)),
	)
	logger.Debug("get contract by ID started")

	res, err = tx(ctx, s.pool, func(pgTx pgx.Tx) (*model.GetContractByIDResponse, error) {
		txLogger := logger.With(
			slog.String("layer", "transaction"),
		)
		txLogger.Debug("transaction execution started")

		contract, err := s.useCase.GetContractByID(ctx, pgTx, s.repo, request)
		if err != nil {
			txLogger.Error("get contract by ID useCase failed", slog.Any("error", err))
			return nil, fmt.Errorf("useCase.GetContractByID: %w", err)
		}

		resp := &model.GetContractByIDResponse{
			ContractID:     contract.ID,
			ContractNumber: contract.Number,
			ContractDate:   contract.Date,
			ContractStatus: contract.Status,
		}

		txLogger.Debug("transaction get contract by ID completed", slog.String("contract_id", strconv.Itoa(resp.ContractID)))
		return resp, nil
	})

	if err != nil {
		logger.Error("get contract by ID failed", slog.Any("error", err))
		return nil, err
	}

	logger.Debug("get contract by ID completed", slog.String("contract_id", strconv.Itoa(res.ContractID)))
	return res, nil
}
