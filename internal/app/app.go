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
func New(pool *pgxpool.Pool) *fiber.App {
	validate := validator.New(validator.WithRequiredStructEnabled())

	infoRepo := storage.NewInfoRepository(pool)
	healthcheckUC := usecase.NewHealthcheckUseCase(infoRepo)
	healthcheckService := service.NewHealthcheckService(healthcheckUC)

	contractRepo := storage.NewContractRepository(pool)
	contractUC := usecase.NewContractUseCase()
	contractSvc := service.NewContractService(pool, contractRepo, contractUC)

	offeringRepo := storage.NewOfferingRepository(pool)
	linkRepo := storage.NewContractOfferingRepository(pool)
	offeringUC := usecase.NewOfferingUseCase()
	offeringSvc := service.NewOfferingService(pool, contractRepo, offeringRepo, linkRepo, offeringUC)

	httpSrv := httpserver.NewServer(validate, healthcheckService, contractSvc, offeringSvc)
	fiberApp := fiber.New()
	httpSrv.SetupRoutes(fiberApp)
	return fiberApp
}
