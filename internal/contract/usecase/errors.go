package usecase

import "errors"

var (
	ErrContractNotFound = errors.New("contract not found")
	ErrServiceNotFound  = errors.New("service not found")
	ErrInvalidRequest   = errors.New("invalid request")
)
