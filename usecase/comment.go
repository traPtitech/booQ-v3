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
	CommentRepository domain.CommentRepository
}

func NewCommentUsecase(commentRepository domain.CommentRepository) CommentUsecase {
	return &commentUsecase{
		CommentRepository: commentRepository,
	}
}

func (u *commentUsecase) CreateComment(itemId int, userId string, text string) (*domain.Comment, error) {

	// そもそもアイテムがない場合は別のエラーを出す
	// そのためには、domain.ItemRepository への依存が必要？

	// 先頭と末尾の空白を削除する
	// 空白だけのコメントを投稿できないようにできる。これは必要かどうか？

	if text == "" {
		return nil, domain.ErrCommentTextEmpty
	}

	comment := &domain.Comment{
		ItemID: itemId,
		UserID: userId,
		Text:   text,
	}

	err := u.CommentRepository.Create(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
