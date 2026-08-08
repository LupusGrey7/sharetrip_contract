package http

import (
	"job4j/sharetrip-contract/internal/contract/service"

	"github.com/go-playground/validator/v10"
)

// Server — HTTP-адаптер (как internal/api.Server в прошлых проектах лида).
// Routes и handler-методы живут в этом же пакете http.
// Сервисы — через interface, чтобы handler-тесты могли подставлять stub.
type Server struct {
	Validator          *validator.Validate
	HealthcheckService service.BaseHealthcheck
	ContractService    service.Contract
	OfferingService    service.Offering
	CompanyService     service.Company
}

func NewServer(
	validator *validator.Validate,
	healthcheckService service.BaseHealthcheck,
	contractService service.Contract,
	offeringService service.Offering,
	companyService service.Company,
) *Server {
	return &Server{
		Validator:          validator,
		HealthcheckService: healthcheckService,
		ContractService:    contractService,
		OfferingService:    offeringService,
		CompanyService:     companyService,
	}
}
