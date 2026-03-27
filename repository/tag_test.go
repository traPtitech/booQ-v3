package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

func TestTagRepository_GetByItemID(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB) (int, []*domain.Tag)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) (int, []*domain.Tag) {
				models := []*tag{
					{ItemID: 1, Name: "go"},
					{ItemID: 1, Name: "book"},
					{ItemID: 2, Name: "equipment"},
				}
				for _, model := range models {
					if err := db.Create(model).Error; err != nil {
						t.Fatalf("failed to create test tag: %v", err)
					}
				}

				return 1, []*domain.Tag{
					{ItemID: 1, Name: "book"},
					{ItemID: 1, Name: "go"},
				}
			},
			wantErr: false,
		},
		{
			name: "success: empty result",
			setup: func(t *testing.T, db *gorm.DB) (int, []*domain.Tag) {
				return 9999, []*domain.Tag{}
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewTagRepository(db)
			itemID, expected := tc.setup(t, db)

			tags, err := repo.GetByItemID(itemID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.ElementsMatch(t, expected, tags)
		})
	}
}

func TestTagRepository_GetByItemIDs(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB) ([]int, map[int][]*domain.Tag)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) ([]int, map[int][]*domain.Tag) {
				models := []*tag{
					{ItemID: 1, Name: "go"},
					{ItemID: 1, Name: "book"},
					{ItemID: 2, Name: "equipment"},
					{ItemID: 3, Name: "unused"},
				}
				for _, model := range models {
					if err := db.Create(model).Error; err != nil {
						t.Fatalf("failed to create test tag: %v", err)
					}
				}

				return []int{1, 2, 9999}, map[int][]*domain.Tag{
					1: {
						{ItemID: 1, Name: "book"},
						{ItemID: 1, Name: "go"},
					},
					2: {
						{ItemID: 2, Name: "equipment"},
					},
					9999: {},
				}
			},
			wantErr: false,
		},
		{
			name: "success: empty item IDs",
			setup: func(t *testing.T, db *gorm.DB) ([]int, map[int][]*domain.Tag) {
				return []int{}, map[int][]*domain.Tag{}
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewTagRepository(db)
			itemIDs, expected := tc.setup(t, db)

			tagsByItemID, err := repo.GetByItemIDs(itemIDs)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(expected), len(tagsByItemID))
			for itemID, expectedTags := range expected {
				assert.ElementsMatch(t, expectedTags, tagsByItemID[itemID])
			}
		})
	}
}

func TestTagRepository_ReplaceByItemID(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(t *testing.T, db *gorm.DB) (int, []string)
		expected []*domain.Tag
		wantErr  bool
	}{
		{
			name: "success: replace existing tags",
			setup: func(t *testing.T, db *gorm.DB) (int, []string) {
				initial := []*tag{
					{ItemID: 1, Name: "old1"},
					{ItemID: 1, Name: "old2"},
				}
				for _, model := range initial {
					db.Create(model)
				}
				return 1, []string{"go", "book"}
			},
			expected: []*domain.Tag{
				{ItemID: 1, Name: "book"},
				{ItemID: 1, Name: "go"},
			},
			wantErr: false,
		},
		{
			name: "success: clear tags",
			setup: func(t *testing.T, db *gorm.DB) (int, []string) {
				initial := []*tag{
					{ItemID: 2, Name: "old"},
				}
				for _, model := range initial {
					db.Create(model)
				}
				return 2, []string{}
			},
			expected: []*domain.Tag{},
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewTagRepository(db)
			itemID, names := tc.setup(t, db)

			err := repo.ReplaceByItemID(itemID, names)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			tags, err := repo.GetByItemID(itemID)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.expected, tags)
		})
	}
}
