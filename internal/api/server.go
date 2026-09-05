package api

import (
	"job4j/sharetrip-contract/internal/contract/service"

	"github.com/go-playground/validator/v10"
)

// Server — HTTP-адаптер (пакет internal/api, как у лида).
// Routes и handler-методы живут в этом же пакете.
// Сервисы — через interface, чтобы handler-тесты могли подставлять stub.
type Server struct {
	Validator          *validator.Validate
	HealthcheckService service.BaseHealthcheck
	ContractService    service.Contract
	CompanyService     service.Company
}

func NewServer(
	validator *validator.Validate,
	healthcheckService service.BaseHealthcheck,
	contractService service.Contract,
	companyService service.Company,
) *Server {
	return &Server{
		Validator:          validator,
		HealthcheckService: healthcheckService,
		ContractService:    contractService,
		CompanyService:     companyService,
	}
}
