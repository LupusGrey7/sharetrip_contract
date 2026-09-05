package api

import (
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var openAPIHostLine = regexp.MustCompile(`(?m)^host:\s*".*"$`)

// candidate paths: make run from repo root, or binary started from cmd/, build/, etc.
func openAPISpecCandidates() []string {
	return []string{
		filepath.Join("api", "contract.yaml"),
		filepath.Join("..", "api", "contract.yaml"),
		filepath.Join("..", "..", "api", "contract.yaml"),
	}
}

// GetOpenAPISpec returns the contents of api/contract.yaml (does not open Swagger Editor itself).
// Editor: download/copy the response or import the file locally — see cheatsheets/swagger-ui-cheatsheet.md.
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
		data = patchOpenAPIHost(data, ctx.Get("Host"))
		ctx.Set(fiber.HeaderContentType, "application/yaml; charset=utf-8")
		return ctx.Status(fiber.StatusOK).Send(data)
	}

	logger.Error("openapi spec not found", slog.Any("error", lastErr))
	return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "api/contract.yaml not found; run the service from repository root (make run)",
	})
}

// patchOpenAPIHost rewrites swagger 2.0 host so Try it out hits the running app (ctx.Host()).
func patchOpenAPIHost(yaml []byte, host string) []byte {
	if host == "" || !openAPIHostLine.Match(yaml) {
		return yaml
	}
	return openAPIHostLine.ReplaceAll(yaml, []byte(`host: "`+host+`"`))
}
