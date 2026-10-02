package api

import (
	"log/slog"
	"os"

	"job4j/sharetrip-contract/gen"

	"github.com/gofiber/fiber/v2"
)

const (
	GroupPrefixV2   = "/api/v2"
	HealthcheckPath = "/healthcheck"
	ContractPath    = "/contracts"
	CompaniesPath   = "/companies"
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
		{Method: "GET", Path: GroupPrefixV2 + ContractPath + "/openapi", Note: "OpenAPI yaml (Swagger UI :8086)"},
		{Method: "POST", Path: GroupPrefixV2 + ContractPath, Note: "create contract"},
		{Method: "GET", Path: GroupPrefixV2 + ContractPath + "/active", Note: "get active contract by companyId query"},
		{Method: "GET", Path: GroupPrefixV2 + CompaniesPath + "/{companyId}/services/{serviceCode}/availability", Note: "check service availability for company"},
		{Method: "GET", Path: GroupPrefixV2 + ContractPath + "/{contractId}", Note: "get by id (stub 501)"},
		{Method: "PATCH", Path: GroupPrefixV2 + ContractPath + "/{contractId}", Note: "update status (stub 501)"},
		{Method: "POST", Path: GroupPrefixV2 + ContractPath + "/{contractId}/signature", Note: "sign (stub 501)"},
		{Method: "PUT", Path: GroupPrefixV2 + "/services", Note: "upsert services (stub 501)"},
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

// SetupRoutes — healthcheck вручную; business paths из gen.RegisterHandlers.
// *Server реализует рабочие handlers с generated query/path params.
// oapiServer добавляет только временные stubs.
func (s *Server) SetupRoutes(app *fiber.App) {
	app.Get(HealthcheckPath, s.Healthcheck)
	gen.RegisterHandlersWithOptions(app, oapiServer{Server: s}, gen.FiberServerOptions{
		BaseURL: GroupPrefixV2,
	})
}
