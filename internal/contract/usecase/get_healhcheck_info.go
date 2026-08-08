package usecase

import (
	"context"
	"fmt"
	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/jackc/pgx/v5"
)

func (u *InfoUseCase) GetHealthcheckInfo(ctx context.Context, tx pgx.Tx, repo storage.BaseTxInfoRepository) (*model.GetHealthcheckInfoResponse, error) {
	info, err := repo.GetHealthcheckInfoTx(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("GetHealthcheckInfo: %w", err)
	}
	return info, nil
}
