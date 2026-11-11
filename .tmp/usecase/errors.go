package usecase

import "errors"

var (
	ErrAlreadyLiked       = errors.New("already liked")
	ErrNotLiked           = errors.New("not liked")
	ErrInvalidSearchQuery = errors.New("offset requires limit")

	// 未決事項や説明だけに留めた処理を、成功したように返さないための目印。
	ErrDesignPending = errors.New("design decision pending; see TODO")
	ErrConceptOnly   = errors.New("concept only: processing omitted")
)
