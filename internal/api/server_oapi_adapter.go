package api

import (
	"job4j/sharetrip-contract/gen"

	openapitypes "github.com/oapi-codegen/runtime/types"

	"github.com/gofiber/fiber/v2"
)

// oapiServer — thin adapter for gen.ServerInterface / RegisterHandlers.
// The server directly implements operations whose signatures match the generated interface.
// Only temporary stubs remain here for operations not yet implemented.
type oapiServer struct {
	*Server
}

var _ gen.ServerInterface = oapiServer{}

// Stubs for OpenAPI operations not implemented yet.

func (a oapiServer) GetContractByID(c *fiber.Ctx, contractId openapitypes.UUID) error {
	_ = contractId
	return fiber.NewError(fiber.StatusNotImplemented, "getContractByID not implemented")
}

func (a oapiServer) UpdateContractStatusById(c *fiber.Ctx, contractId openapitypes.UUID) error {
	_ = contractId
	return fiber.NewError(fiber.StatusNotImplemented, "updateContractStatusById not implemented")
}

func (a oapiServer) SignContract(c *fiber.Ctx, contractId openapitypes.UUID) error {
	_ = contractId
	return fiber.NewError(fiber.StatusNotImplemented, "signContract not implemented")
}

func (a oapiServer) UpsertContractServices(c *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotImplemented, "upsertContractServices not implemented")
}
