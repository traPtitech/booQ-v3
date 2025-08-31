package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOwnershipTableName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "ownerships", (&Ownership{}).TableName())
}

func TestGetOwnership(t *testing.T) {
	PrepareTestDatabase()

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		res, err := GetOwnership(-1)
		assert.Error(err)
		assert.Empty(res)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		res, err := GetOwnership(1)
		assert.NoError(err)
		assert.NotEmpty(res)
		assert.Equal(res.UserID, "s9")
		assert.Equal(res.ItemID, 1)
		assert.Equal(res.Rentalable, true)
		assert.Equal(res.Memo, "memo1")
		assert.Equal(4, len(res.Transaction))
	})
}

func TestCreateOwnership(t *testing.T) {
	PrepareTestDatabase()

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)

		payload := OwnershipPayload{
			ItemID:     1,
			UserID:     "new_user",
			Rentalable: true,
			Memo:       "created by test",
		}
		created, err := CreateOwnership(payload)
		assert.NoError(err)
		assert.NotZero(created.ID)
		assert.Equal(payload.ItemID, created.ItemID)
		assert.Equal(payload.UserID, created.UserID)
		assert.Equal(payload.Rentalable, created.Rentalable)
		assert.Equal(payload.Memo, created.Memo)

		fetched, err := GetOwnership(created.ID)
		assert.NoError(err)
		assert.Equal(created.ID, fetched.ID)
		assert.Equal(created.UserID, fetched.UserID)
	})
}

func TestUpdateOwnership(t *testing.T) {
	PrepareTestDatabase()

	t.Run("failure: not found", func(t *testing.T) {
		assert := assert.New(t)
		_, err := UpdateOwnership(-1, OwnershipPayload{})
		assert.ErrorIs(err, ErrNotFound)
	})

	t.Run("failure: unauthorized", func(t *testing.T) {
		assert := assert.New(t)
		_, err := UpdateOwnership(1, OwnershipPayload{
			ItemID:     1,
			UserID:     "different_user",
			Rentalable: false,
			Memo:       "new_memo",
		})
		assert.ErrorIs(err, ErrUnauthorized)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)
		updated, err := UpdateOwnership(1, OwnershipPayload{
			ItemID:     1,
			UserID:     "s9",
			Rentalable: false,
			Memo:       "updated",
		})
		assert.NoError(err)
		assert.Equal("s9", updated.UserID)
		assert.Equal(false, updated.Rentalable)
		assert.Equal("updated", updated.Memo)

		fetched, err := GetOwnership(1)
		assert.NoError(err)
		assert.Equal("updated", fetched.Memo)
		assert.Equal(false, fetched.Rentalable)
	})
}

func TestDeleteOwnership(t *testing.T) {
	PrepareTestDatabase()

	t.Run("failure: not found", func(t *testing.T) {
		assert := assert.New(t)
		err := DeleteOwnership(-1, "someone")
		assert.ErrorIs(err, ErrNotFound)
	})

	t.Run("failure: unauthorized", func(t *testing.T) {
		assert := assert.New(t)
		err := DeleteOwnership(1, "different_user")
		assert.ErrorIs(err, ErrUnauthorized)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)
		err := DeleteOwnership(1, "s9")
		assert.NoError(err)

		res, err := GetOwnership(1)
		assert.ErrorIs(err, ErrNotFound)
		assert.Empty(res)
	})
}
