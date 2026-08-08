package usecase

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

type BaseCompanyUseCase interface {
	GetAvailableOfferingByCompanyID(
		ctx context.Context,
		tx pgx.Tx,
		offeringRepo storage.BaseTxOfferingRepository,
		companyRepo storage.BaseCompanyRepository,
		req *model.GetAvailableOfferingByCompanyIDRequest,
	) (*model.GetAvailableResultResponse, error)
}

type CompanyUseCase struct{}

func NewCompanyUseCase() *CompanyUseCase {
	return &CompanyUseCase{}
}
