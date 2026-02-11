package usecase

import (
	"github.com/traPtitech/booQ-v3/domain"
)

type CommentUsecase interface {
	CreateComment(
		itemID int,
		userID string,
		text string,
	) (
		*domain.Comment, error,
	)
}

type commentUsecase struct {
	CommentRepo domain.CommentRepository
	ItemRepo    domain.ItemRepository
}

func NewCommentUsecase(commentRepo domain.CommentRepository, ItemRepo domain.ItemRepository) CommentUsecase {
	return &commentUsecase{
		CommentRepo: commentRepo,
		ItemRepo:    ItemRepo,
	}
}

func (u *commentUsecase) CreateComment(itemId int, userId string, text string) (*domain.Comment, error) {

	_, err := u.ItemRepo.GetByID(itemId)

	if err != nil {
		return nil, err
	}

	if text == "" {
		return nil, domain.ErrCommentTextEmpty
	}

	comment := &domain.Comment{
		ItemID: itemId,
		UserID: userId,
		Text:   text,
	}

	err = u.CommentRepo.Create(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
