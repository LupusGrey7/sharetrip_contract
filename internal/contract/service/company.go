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
	GetAvailableServiceByCompanyID(ctx context.Context, req *model.GetAvailableServiceByCompanyIDRequest) (*model.GetAvailableResultResponse, error)
}

type CompanyService struct {
	pool        *pgxpool.Pool
	companyRepo storage.BaseCompanyRepository
	useCase     usecase.BaseCompanyUseCase
}

func NewCompanyService(useCase usecase.BaseCompanyUseCase) *CompanyService {
	return &CompanyService{useCase: useCase}
}
