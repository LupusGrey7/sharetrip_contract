package http

import (
	"github.com/gofiber/fiber/v2"
)

const (
	GroupPrefixV2   = "/api/v2"
	HealthcheckPath = "/healthcheck"
	ContractPath    = "/contracts"
	ServicesPath    = "/services"

	contractGetByIdPath = "/:contractId"
)

func (s *Server) SetupRoutes(app *fiber.App) {
	app.Get(HealthcheckPath, s.Healthcheck)

	v2 := app.Group(GroupPrefixV2)
	v2.Put(ServicesPath, s.UpsertContractServices)

	contracts := v2.Group(ContractPath)
	contracts.Post("/", s.CreateContract)
	contracts.Get(contractGetByIdPath, s.GetContractByID)
}
