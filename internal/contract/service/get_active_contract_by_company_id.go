package service

import (
	"context"
	"fmt"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
)

func (s *ContractService) GetActiveContractByCompanyID(
	ctx context.Context,
	input *domain.GetActiveContractByCompanyIDInput,
) (*domain.ContractOutput, error) {
	res, err := tx(ctx, s.pool, func(pgTx pgx.Tx) (*domain.ContractOutput, error) {
		return s.useCase.GetActiveContractByCompanyID(ctx, pgTx, s.repo, input)
	})
	if err != nil {
		return nil, fmt.Errorf("GetActiveContractByCompanyID: %w", err)
	}
	return res, nil
}
