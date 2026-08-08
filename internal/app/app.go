package app

import (
	"job4j/sharetrip-contract/internal/contract/service"
	"job4j/sharetrip-contract/internal/contract/usecase"
	httpserver "job4j/sharetrip-contract/internal/http"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// New is the composition root: repo -> usecase -> service -> http.Server -> Fiber.
//
// Why not only http.NewServer?
//
//	http.NewServer — HTTP adapter; needs already-built services.
//	app.New        — builds the full dependency chain from pool.
//
// If wiring stays in main, package app is pointless.
func New(pool *pgxpool.Pool) *fiber.App {
	validate := validator.New(validator.WithRequiredStructEnabled())

	infoRepo := storage.NewInfoRepository(pool)
	healthcheckUC := usecase.NewHealthcheckUseCase(infoRepo)
	healthcheckService := service.NewHealthcheckService(healthcheckUC)

	contractRepo := storage.NewContractRepository(pool)
	contractUC := usecase.NewContractUseCase()
	contractSvc := service.NewContractService(pool, contractRepo, contractUC)

	offeringRepo := storage.NewOfferingRepository(pool)
	offeringUC := usecase.NewOfferingUseCase()
	offeringSvc := service.NewOfferingService(pool, offeringRepo, offeringUC)

	httpSrv := httpserver.NewServer(validate, healthcheckService, contractSvc, offeringSvc)
	fiberApp := fiber.New()
	httpSrv.SetupRoutes(fiberApp)
	return fiberApp
}
