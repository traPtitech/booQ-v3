package model

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrUpdateNotAllowed = errors.New("some fields cannot be updated")
	ErrUnauthorized     = errors.New("unauthorized")
)
