package api

import (
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

// Swagger 2.0 leftover (no longer in contract.yaml).
var openAPIHostLine = regexp.MustCompile(`(?m)^host:\s*".*"$`)

// OpenAPI 3.x: first servers[].url with localhost (local Try it out).
var openAPILocalServerURL = regexp.MustCompile(`(?m)^([ \t]*(?:-[ \t]*)?url:[ \t]*)https?://localhost:\d+/api/v2[ \t]*$`)

// Prefer the bundled document produced by make generate: Swagger UI cannot fetch
// repository-local $ref files through this single HTTP endpoint.
// Fall back to the root document for development before the first generation.
func openAPISpecCandidates() []string {
	return []string{
		filepath.Join("build", "openapi.bundle.yaml"),
		filepath.Join("api", "contract.yaml"),
		filepath.Join("..", "build", "openapi.bundle.yaml"),
		filepath.Join("..", "api", "contract.yaml"),
		filepath.Join("..", "..", "build", "openapi.bundle.yaml"),
		filepath.Join("..", "..", "api", "contract.yaml"),
	}
}

// GetOpenAPISpec returns the bundled OpenAPI document produced from api/contract.yaml.
// Editor: download/copy the response or import the source root locally.
func (s *Server) GetOpenAPISpec(ctx *fiber.Ctx) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "http"),
		slog.String("handler", "GetOpenAPISpec"),
	)

	var lastErr error
	for _, path := range openAPISpecCandidates() {
		data, err := os.ReadFile(path)
		if err != nil {
			lastErr = err
			continue
		}
		logger.Debug("openapi spec served", slog.String("path", path))
		data = patchOpenAPIServer(data, ctx.Get("Host"))
		ctx.Set(fiber.HeaderContentType, "application/yaml; charset=utf-8")
		return ctx.Status(fiber.StatusOK).Send(data)
	}

	logger.Error("openapi spec not found", slog.Any("error", lastErr))
	return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "api/contract.yaml not found; run the service from repository root (make run)",
	})
}

// patchOpenAPIServer rewrites local server URL (OAS3) or host (Swagger 2) so Try it out hits the running app.
func patchOpenAPIServer(yaml []byte, host string) []byte {
	if host == "" {
		return yaml
	}
	if openAPILocalServerURL.Match(yaml) {
		return openAPILocalServerURL.ReplaceAll(yaml, []byte("${1}http://"+host+"/api/v2"))
	}
	if openAPIHostLine.Match(yaml) {
		return openAPIHostLine.ReplaceAll(yaml, []byte(`host: "`+host+`"`))
	}
	return yaml
}
