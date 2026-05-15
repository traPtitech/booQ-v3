package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
)

type stubCommentRepository struct {
	createFn func(comment *domain.Comment) (*domain.Comment, error)
	created  []*domain.Comment
}

func (s *stubCommentRepository) Create(comment *domain.Comment) (*domain.Comment, error) {
	s.created = append(s.created, comment)
	if s.createFn != nil {
		return s.createFn(comment)
	}
	return comment, nil
}

type stubItemRepository struct {
	getByIDFn    func(id int) (*domain.Item, error)
	getByIDCalls []int
}

func (s *stubItemRepository) GetByID(id int) (*domain.Item, error) {
	s.getByIDCalls = append(s.getByIDCalls, id)
	if s.getByIDFn != nil {
		return s.getByIDFn(id)
	}
	return &domain.Item{ID: id}, nil
}

func (s *stubItemRepository) Search(query domain.ItemSearchQuery) ([]*domain.Item, error) {
	return nil, errors.New("not implemented")
}

func (s *stubItemRepository) Create(item *domain.Item) error {
	return errors.New("not implemented")
}

func (s *stubItemRepository) Update(item *domain.Item) error {
	return errors.New("not implemented")
}

func (s *stubItemRepository) Delete(id int) error {
	return errors.New("not implemented")
}

func TestCommentUsecase_CreateComment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		itemRepo := &stubItemRepository{}
		commentRepo := &stubCommentRepository{}
		u := NewCommentUsecase(commentRepo, itemRepo)

		comment, err := u.CreateComment(1, "user-1", "hello")

		assert.NoError(t, err)
		assert.Equal(t, []int{1}, itemRepo.getByIDCalls)
		assert.Len(t, commentRepo.created, 1)
		assert.Equal(t, 1, comment.ItemID)
		assert.Equal(t, "user-1", comment.UserID)
		assert.Equal(t, "hello", comment.Text)
	})

	t.Run("item not found", func(t *testing.T) {
		itemRepo := &stubItemRepository{
			getByIDFn: func(id int) (*domain.Item, error) {
				return nil, domain.ErrNotFound
			},
		}
		commentRepo := &stubCommentRepository{}
		u := NewCommentUsecase(commentRepo, itemRepo)

		comment, err := u.CreateComment(99, "user-1", "hello")

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Equal(t, []int{99}, itemRepo.getByIDCalls)
		assert.Empty(t, commentRepo.created)
	})

	t.Run("empty comment text", func(t *testing.T) {
		itemRepo := &stubItemRepository{}
		commentRepo := &stubCommentRepository{}
		u := NewCommentUsecase(commentRepo, itemRepo)

		comment, err := u.CreateComment(1, "user-1", "")

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, domain.ErrCommentTextEmpty)
		assert.Equal(t, []int{1}, itemRepo.getByIDCalls)
		assert.Empty(t, commentRepo.created)
	})

	t.Run("comment repository error", func(t *testing.T) {
		repoErr := errors.New("insert failed")
		itemRepo := &stubItemRepository{}
		commentRepo := &stubCommentRepository{
			createFn: func(comment *domain.Comment) (*domain.Comment, error) {
				return nil, repoErr
			},
		}
		u := NewCommentUsecase(commentRepo, itemRepo)

		comment, err := u.CreateComment(1, "user-1", "hello")

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, []int{1}, itemRepo.getByIDCalls)
		assert.Len(t, commentRepo.created, 1)
	})
}
