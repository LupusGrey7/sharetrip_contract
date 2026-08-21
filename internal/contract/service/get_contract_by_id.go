package service

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
)

func (s *ContractService) GetContractByID(
	ctx context.Context,
	input *domain.GetContractByIDInput,
) (*domain.ContractOutput, error) {
	res, err := tx(ctx, s.pool, func(pgTx pgx.Tx) (*domain.ContractOutput, error) {
		return s.useCase.GetContractByID(ctx, pgTx, s.repo, input)
	})
	if err != nil {
		return nil, fmt.Errorf("GetContractByID: %w", err)
	}
	return res, nil
}
