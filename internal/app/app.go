package app

import (
	api "job4j/sharetrip-contract/internal/api"
	"job4j/sharetrip-contract/internal/contract/service"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
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
	linkRepo := storage.NewContractOfferingRepository(pool)
	offeringUC := usecase.NewOfferingUseCase()
	offeringSvc := service.NewOfferingService(pool, contractRepo, offeringRepo, linkRepo, offeringUC)

	companyRepo := storage.NewCompanyRepository(pool)
	companyUC := usecase.NewCompanyUseCase()
	companySvc := service.NewCompanyService(pool, offeringRepo, companyRepo, companyUC)

	httpSrv := api.NewServer(validate, healthcheckService, contractSvc, offeringSvc, companySvc)
	fiberApp := fiber.New()
	httpSrv.SetupRoutes(fiberApp)
	return fiberApp
}
