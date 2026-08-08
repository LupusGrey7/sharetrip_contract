package http

import (
	"github.com/gofiber/fiber/v2"
)

const (
	GroupPrefixV2                    = "/api/v2"
	HealthcheckPath                  = "/healthcheck"
	ContractPath                     = "/contracts"
	contractGetByIdPath              = "/:contractId"
	contractGetActivePath            = "/active"
	contractGetActiveByCompanyIdPath = "/active/:companyId"
)

func (s *Server) SetupRoutes(app *fiber.App) {
	app.Get(HealthcheckPath, s.Healthcheck)

	v2 := app.Group(GroupPrefixV2)
	contracts := v2.Group(ContractPath)
	contracts.Get(contractGetByIdPath, s.GetContractByID)
	_ = contractGetActivePath
	_ = contractGetActiveByCompanyIdPath
}
