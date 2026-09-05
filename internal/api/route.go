package api

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

const (
	GroupPrefixV2   = "/api/v2"
	HealthcheckPath = "/healthcheck"
	ContractPath    = "/contracts"
	CompaniesPath   = "/companies"
	contractOpenAPIPath            = "/openapi"
	contractGetActivePath          = "/active"
	companyServiceAvailabilityPath = "/:companyId/services/:serviceCode/availability"
)

// RouteInfo — method + path for logging at startup.
type RouteInfo struct {
	Method string
	Path   string
	Note   string
}

// RegisteredRoutes — list of HTTP API that SetupRoutes starts.
func RegisteredRoutes() []RouteInfo {
	return []RouteInfo{
		{Method: "GET", Path: HealthcheckPath, Note: "liveness / DB ping"},
		{Method: "GET", Path: GroupPrefixV2 + ContractPath + contractOpenAPIPath, Note: "OpenAPI yaml (Swagger UI :8086)"},
		{Method: "POST", Path: GroupPrefixV2 + ContractPath + "/", Note: "create contract"},
		{Method: "GET", Path: GroupPrefixV2 + ContractPath + "/active", Note: "get active contract by companyId query"},
		{Method: "GET", Path: GroupPrefixV2 + CompaniesPath + "/{companyId}/services/{serviceCode}/availability", Note: "check service availability for company"},
	}
}

// LogRegisteredRoutes — logging available endpoints.
func LogRegisteredRoutes(addr string) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(slog.String("layer", "http"))
	logger.Info("available endpoints", slog.String("listen", addr))
	for _, r := range RegisteredRoutes() {
		if r.Note != "" {
			logger.Info("route", slog.String("method", r.Method), slog.String("url", r.Path), slog.String("note", r.Note))
		} else {
			logger.Info("route", slog.String("method", r.Method), slog.String("url", r.Path))
		}
	}
}

// SetupRoutes — setup routes for the server.
func (s *Server) SetupRoutes(app *fiber.App) {
	app.Get(HealthcheckPath, s.Healthcheck)

	v2 := app.Group(GroupPrefixV2)

	contracts := v2.Group(ContractPath)
	contracts.Get(contractOpenAPIPath, s.GetOpenAPISpec)
	contracts.Post("/", s.CreateContract)
	contracts.Get(contractGetActivePath, s.GetActiveContractByCompanyID)

	companies := v2.Group(CompaniesPath)
	companies.Get(companyServiceAvailabilityPath, s.GetAvailableOfferingByCompanyID)
}
