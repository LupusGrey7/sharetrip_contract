package http

import (
	"job4j/sharetrip-contract/internal/contract/service"

	"github.com/go-playground/validator/v10"
)

// Server — HTTP-адаптер (как internal/api.Server в прошлых проектах лида).
// Routes и методы Get*/Healthcheck живут в этом же пакете http.
type Server struct {
	Validator          *validator.Validate
	HealthcheckService *service.HealthcheckService
	ContractService    *service.ContractService
	OfferingService    *service.OfferingService
}

func NewServer(
	validator *validator.Validate,
	healthcheckService *service.HealthcheckService,
	contractService *service.ContractService,
	offeringService *service.OfferingService,
) *Server {
	return &Server{
		Validator:          validator,
		HealthcheckService: healthcheckService,
		ContractService:    contractService,
		OfferingService:    offeringService,
	}
}
