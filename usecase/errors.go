package usecase

import "errors"

var (
	ErrInvalidSearchQuery = errors.New("invalid search query")
	ErrUpdateNotAllowed   = errors.New("some fields cannot be updated")
	ErrForbidden          = errors.New("you cannot perform this action")
)
