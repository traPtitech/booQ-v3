package domain

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrInvalidOwner         = errors.New("invalid owner")
	ErrInvalidQuantity      = errors.New("invalid quantity")
	ErrInvalidPolicy        = errors.New("invalid policy")
	ErrLendingDisabled      = errors.New("lending disabled")
	ErrInsufficientQuantity = errors.New("insufficient quantity")
	ErrInvalidState         = errors.New("invalid state")
	ErrInvalidDueDate       = errors.New("due date must be in the future")
	ErrOwnershipInUse       = errors.New("ownership has active loans or requests")
	ErrCommentTextEmpty     = errors.New("comment text cannot be empty")
	ErrFileTooLarge         = errors.New("file exceeds 3 MB")
	ErrInvalidFileType      = errors.New("only JPEG and PNG are supported")
)
