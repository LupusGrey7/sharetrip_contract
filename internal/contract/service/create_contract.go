package service

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
)

func (s *ContractService) CreateContract(
	ctx context.Context,
	input *domain.CreateContractInput,
) (*domain.ContractOutput, error) {
	// pool == nil — only unit tests without DB (direct use case call).
	if s.pool == nil {
		contract, err := s.useCase.CreateContract(ctx, nil, s.repo, input)
		if err != nil {
			return nil, fmt.Errorf("CreateContract: %w", err)
		}
		return contract, nil
	}

	res, err := tx(ctx, s.pool, func(pgTx pgx.Tx) (*domain.ContractOutput, error) {
		return s.useCase.CreateContract(ctx, pgTx, s.repo, input)
	})
	if err != nil {
		return nil, fmt.Errorf("CreateContract: %w", err)
	}
	return res, nil
}
