package service

import (
	"context"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Company — сценарии availability по компании (методы — в отдельных файлах).
type Company interface {
	GetAvailableOfferingByCompanyID(ctx context.Context, req *model.GetAvailableOfferingByCompanyIDRequest) (*model.GetAvailableResultResponse, error)
}

type CompanyService struct {
	pool         *pgxpool.Pool
	offeringRepo storage.BaseTxOfferingRepository
	companyRepo  storage.BaseCompanyRepository
	useCase      usecase.BaseCompanyUseCase
}

func NewCompanyService(
	pool *pgxpool.Pool,
	offeringRepo storage.BaseTxOfferingRepository,
	companyRepo storage.BaseCompanyRepository,
	useCase usecase.BaseCompanyUseCase,
) *CompanyService {
	return &CompanyService{
		pool:         pool,
		offeringRepo: offeringRepo,
		companyRepo:  companyRepo,
		useCase:      useCase,
	}
}
