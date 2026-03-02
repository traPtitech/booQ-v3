package domain

import "errors"

var (
	ErrItemNotFound    = errors.New("item not found")
	ErrFileNotFound    = errors.New("file not found")
	ErrInvalidFileType = errors.New("invalid file type: only JPEG and PNG are allowed")
	ErrFileTooLarge    = errors.New("file too large: max size is 3MB")
)
