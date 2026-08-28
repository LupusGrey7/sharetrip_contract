package app

import (
	"context"

	api "job4j/sharetrip-contract/internal/api"
	"job4j/sharetrip-contract/internal/contract/service"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"job4j/sharetrip-contract/configs"

	"job4j/sharetrip-contract/internal/observability/tracing"
)

// New is the composition root: repo -> usecase -> service -> api.Server -> Fiber.
func New(pool *pgxpool.Pool) *fiber.App {
	validate := validator.New(validator.WithRequiredStructEnabled())

	infoRepo := storage.NewInfoRepository(pool)
	healthcheckUC := usecase.NewHealthcheckUseCase(infoRepo)
	healthcheckService := service.NewHealthcheckService(healthcheckUC)

	contractRepo := storage.NewContractRepository(pool)
	contractUC := usecase.NewContractUseCase()
	contractSvc := service.NewContractService(pool, contractRepo, contractUC)

	offeringRepo := storage.NewOfferingRepository(pool)
	companyRepo := storage.NewCompanyRepository(pool)
	companyUC := usecase.NewCompanyUseCase()
	companySvc := service.NewCompanyService(pool, offeringRepo, companyRepo, companyUC)

	httpSrv := api.NewServer(validate, healthcheckService, contractSvc, companySvc)
	fiberApp := fiber.New()
	httpSrv.SetupRoutes(fiberApp)
	return fiberApp
}

func InitTracing(ctx context.Context) (*tracing.TracerProvider, error) {
	return tracing.NewProvider(ctx, tracing.Config{
		ServiceName:    configs.Env("OTEL_SERVICE_NAME", "share-trip"),
		ServiceVersion: configs.Env("OTEL_SERVICE_VERSION", "1.0.0"),
		Environment:    configs.Env("OTEL_ENVIRONMENT", "local"),
		Endpoint:       configs.Env("OTEL_EXPORTER_ENDPOINT", "localhost:4319"),
	})
}
