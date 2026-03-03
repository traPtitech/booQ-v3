package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidFileType = errors.New("invalid file type: only JPEG and PNG are allowed")
	ErrFileTooLarge    = errors.New("file too large: max size is 3MB")
)
