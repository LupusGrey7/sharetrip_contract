package storage

import "errors"

const (
	errQueryByID          = "error querying by ID: %d: %w"
	errOfferingNotFound   = "offering not found"
	errSelectEntityFailed = "error selecting entity: %w"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrInternal           = errors.New("internal server error")
	ErrInvalidID          = errors.New("invalid ID")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidResponse    = errors.New("invalid response")
	ErrInvalidData        = errors.New("invalid data")
	ErrInvalidTransaction = errors.New("invalid transaction")
	ErrInvalidConnection  = errors.New("invalid connection")
	ErrContractNotFound   = errors.New("contract not found")
	ErrOfferingNotFound   = errors.New("offering not found")
)
