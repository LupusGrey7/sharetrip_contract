package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrInvalidIDParamFormat):
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidIDParamFormat.Error())
	case errors.Is(err, ErrContractNotFound):
		return fiber.NewError(fiber.StatusNotFound, ErrContractNotFound.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}
