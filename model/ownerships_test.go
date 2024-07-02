package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOwnershipTableName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "ownerships", (&Ownership{}).TableName())
}

func TestRegisterOwnership(t *testing.T) {
	PrepareTestDatabase()

	t.Run("failures-invaild-id", func(t *testing.T) {
		assert := assert.New(t)
		ownership, err := RegisterOwnership(Ownership{
			ItemID: -1, UserID: "s9", Rentalable: true, Memo: "test"})
		assert.Error(err)
		assert.Empty(ownership)
	})

	t.Run("failures-setting-owner-to-equipment", func(t *testing.T) {
		assert := assert.New(t)
		ownership, err := RegisterOwnership(Ownership{
			ItemID: 3, UserID: "s9", Rentalable: true, Memo: "test"})
		assert.Error(err)
		assert.Empty(ownership)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)
		ownership, err := RegisterOwnership(Ownership{
			ItemID: 1, UserID: "s9", Rentalable: true, Memo: "test"})
		assert.NoError(err)
		assert.NotEmpty(ownership)
		assert.Equal(ownership.ItemID, 1)
		assert.Equal(ownership.UserID, "s9")
		assert.Equal(ownership.Rentalable, true)
		assert.Equal(ownership.Memo, "test")
	})
}

func TestGetOwnershipByID(t *testing.T) {
	PrepareTestDatabase()
	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		ownership, err := GetOwnershipByID(-1)
		assert.Error(err)
		assert.Empty(ownership)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		ownership, err := GetOwnershipByID(1)
		assert.NoError(err)
		assert.NotEmpty(ownership)
		assert.Equal(ownership.ItemID, 1)
		assert.Equal(ownership.UserID, "s9")
		assert.Equal(ownership.Rentalable, true)
		assert.Equal(ownership.Memo, "memo1")
	})
}

func TestPatchOwnership(t *testing.T) {
	PrepareTestDatabase()

	t.Run("failure", func(t *testing.T) {
		assert := assert.New(t)

		err := PatchOwnership(Ownership{
			GormModel: GormModel{ID: -1},
			ItemID:    1, UserID: "s9", Rentalable: true, Memo: "test"})

		assert.Error(err)
	})

	t.Run("success-1", func(t *testing.T) {
		assert := assert.New(t)

		err := PatchOwnership(Ownership{
			GormModel: GormModel{ID: 1},
			ItemID:    1, UserID: "s9", Rentalable: true, Memo: "test"})
		assert.NoError(err)

		ownership := Ownership{}
		err = db.First(&ownership, 1).Error
		assert.NoError(err)
		assert.Equal(ownership.ItemID, 1)
		assert.Equal(ownership.UserID, "s9")
		assert.Equal(ownership.Rentalable, true)
		assert.Equal(ownership.Memo, "test")
	})

	t.Run("success-2", func(t *testing.T) {
		assert := assert.New(t)

		err := PatchOwnership(Ownership{
			GormModel: GormModel{ID: 1},
			ItemID:    1, UserID: "s9", Rentalable: false, Memo: "test"})
		assert.NoError(err)

		ownership := Ownership{}
		err = db.First(&ownership, 1).Error
		assert.NoError(err)
		assert.Equal(ownership.ItemID, 1)
		assert.Equal(ownership.UserID, "s9")
		assert.Equal(ownership.Rentalable, false)
		assert.Equal(ownership.Memo, "test")
	})
}

func TestDeleteOwnership(t *testing.T) {
	PrepareTestDatabase()

	t.Run("failure", func(t *testing.T) {
		assert := assert.New(t)
		err := DeleteOwnership(-1)
		assert.Error(err)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)
		err := DeleteOwnership(1)
		assert.NoError(err)

		err = db.First(&Ownership{}, 1).Error
		assert.Error(err)

		err = db.First(&Transaction{}, 1).Error
		assert.Error(err)
	})
}
