package service

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/jackc/pgx/v5"
)

func (s *ContractService) CreateContract(
	ctx context.Context,
	request *model.CreateContractRequest,
) (*model.ContractResponse, error) {
	// pool == nil — only unit tests without DB (direct use case call).
	if s.pool == nil {
		contract, err := s.useCase.CreateContract(ctx, nil, s.repo, request)
		if err != nil {
			return nil, fmt.Errorf("CreateContract: %w", err)
		}
		return model.ContractToResponse(contract), nil
	}

	res, err := tx(ctx, s.pool, func(pgTx pgx.Tx) (*model.ContractResponse, error) {
		contract, err := s.useCase.CreateContract(ctx, pgTx, s.repo, request)
		if err != nil {
			return nil, err
		}
		return model.ContractToResponse(contract), nil
	})
	if err != nil {
		return nil, fmt.Errorf("CreateContract: %w", err)
	}
	return res, nil
}
