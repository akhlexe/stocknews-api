package apperrors

import "errors"

// Define sentinel errors (specific error values) that can be checked against

var (
	ErrNotFound           = errors.New("resource not found")
	ErrServiceUnavailable = errors.New("external service unavailable")
	ErrInternal           = errors.New("internal processing error")
	ErrConfiguration      = errors.New("application configuration error")
)


