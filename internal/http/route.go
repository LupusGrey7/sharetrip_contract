package http

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

const (
	GroupPrefixV2   = "/api/v2"
	HealthcheckPath = "/healthcheck"
	ContractPath    = "/contracts"
	ServicesPath    = "/services"
	CompaniesPath   = "/companies"

	contractGetByIdPath              = "/:contractId"
	companyServiceAvailabilityPath   = "/:companyId/services/:serviceCode/availability"
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
		{Method: "POST", Path: GroupPrefixV2 + ContractPath + "/", Note: "create contract"},
		{Method: "GET", Path: GroupPrefixV2 + ContractPath + "/{contractId}", Note: "get contract by id"},
		{Method: "PUT", Path: GroupPrefixV2 + ServicesPath, Note: "upsert contract services (e.g. trip_creation)"},
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
	v2.Put(ServicesPath, s.UpsertContractServices)

	contracts := v2.Group(ContractPath)
	contracts.Post("/", s.CreateContract)
	contracts.Get(contractGetByIdPath, s.GetContractByID)

	companies := v2.Group(CompaniesPath)
	companies.Get(companyServiceAvailabilityPath, s.GetAvailableOfferingByCompanyID)
}
