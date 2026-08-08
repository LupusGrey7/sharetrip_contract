package http

import (
	"errors"

	"job4j/sharetrip-contract/internal/contract/usecase"

	"github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrInvalidIDParamFormat), errors.Is(err, ErrInvalidRequest), errors.Is(err, usecase.ErrInvalidRequest):
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    "COMPANY_VALIDATE_ERROR",
			Message: err.Error(),
		})
	case errors.Is(err, ErrContractNotFound), errors.Is(err, usecase.ErrContractNotFound):
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Code:    "CONTRACT_NOT_FOUND",
			Message: err.Error(),
		})
	case errors.Is(err, usecase.ErrCompanyNotFound):
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Code:    "COMPANY_NOT_FOUND",
			Message: err.Error(),
		})
	case errors.Is(err, usecase.ErrServiceNotFound):
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Code:    "SERVICE_NOT_FOUND",
			Message: err.Error(),
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
}
