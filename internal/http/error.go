package http

import "errors"

const (
	errInvalidRequest       = "invalid request"
	errInvalidIDParamFormat = "invalid id param format"
	errContractNotFound     = "contract not found"
)

var (
	ErrInvalidRequest       = errors.New(errInvalidRequest)
	ErrInvalidIDParamFormat = errors.New(errInvalidIDParamFormat)
	ErrContractNotFound     = errors.New(errContractNotFound)
)
