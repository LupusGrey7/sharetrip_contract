package middleware

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"job4j/sharetrip-contract/internal/observability/logctx"
)

// Correlation puts trace_id and request metadata into UserContext for slog downstream.
// Requires tracing.NewFiberMiddleware() registered before this handler.
func Correlation() fiber.Handler {
	base := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		if traceID != "" && traceID != "00000000000000000000000000000000" {
			c.Set("X-Request-ID", traceID)
		}

		logger := base.With(
			slog.String("trace_id", traceID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
		)
		ctx = logctx.WithLogger(ctx, logger)
		if traceID != "" && traceID != "00000000000000000000000000000000" {
			ctx = logctx.WithRequestID(ctx, traceID)
		}
		c.SetUserContext(ctx)
		return c.Next()
	}
}
