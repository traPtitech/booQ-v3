package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	mock_domain "github.com/traPtitech/booQ-v3/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestCommentUsecase_CreateComment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		itemRepo := mock_domain.NewMockItemRepository(ctrl)
		commentRepo := mock_domain.NewMockCommentRepository(ctrl)

		itemRepo.EXPECT().GetByID(1).Return(&domain.Item{ID: 1}, nil)
		commentRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(c *domain.Comment) (*domain.Comment, error) {
			c.ID = 10
			return c, nil
		})

		u := NewCommentUsecase(commentRepo, itemRepo)
		comment, err := u.CreateComment(1, "user-1", "hello")

		assert.NoError(t, err)
		assert.Equal(t, 1, comment.ItemID)
		assert.Equal(t, "user-1", comment.UserID)
		assert.Equal(t, "hello", comment.Text)
	})

	t.Run("item not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		itemRepo := mock_domain.NewMockItemRepository(ctrl)
		commentRepo := mock_domain.NewMockCommentRepository(ctrl)

		itemRepo.EXPECT().GetByID(99).Return(nil, domain.ErrNotFound)

		u := NewCommentUsecase(commentRepo, itemRepo)
		comment, err := u.CreateComment(99, "user-1", "hello")

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("empty comment text", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		itemRepo := mock_domain.NewMockItemRepository(ctrl)
		commentRepo := mock_domain.NewMockCommentRepository(ctrl)

		itemRepo.EXPECT().GetByID(1).Return(&domain.Item{ID: 1}, nil)

		u := NewCommentUsecase(commentRepo, itemRepo)
		comment, err := u.CreateComment(1, "user-1", "")

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, domain.ErrCommentTextEmpty)
	})

	t.Run("comment repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoErr := errors.New("insert failed")
		itemRepo := mock_domain.NewMockItemRepository(ctrl)
		commentRepo := mock_domain.NewMockCommentRepository(ctrl)

		itemRepo.EXPECT().GetByID(1).Return(&domain.Item{ID: 1}, nil)
		commentRepo.EXPECT().Create(gomock.Any()).Return(nil, repoErr)

		u := NewCommentUsecase(commentRepo, itemRepo)
		comment, err := u.CreateComment(1, "user-1", "hello")

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, repoErr)
	})
}
