package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemTableName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "items", (&Item{}).TableName())
}

func TestGetItems(t *testing.T) {
	PrepareTestDatabase()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		res, err := GetItems(GetItemsBody{Search: "item-id1"})
		assert.NoError(err)
		assert.NotEmpty(res)
		assert.Equal(res[0].Name, "item-id1")
		assert.Equal(res[0].Description, "aaa")
		assert.Equal(res[0].ImgURL, "url")
		assert.Empty(res[0].Book)
		assert.Empty(res[0].Equipment)
	})
}

func TestCreateItems(t *testing.T) {
	PrepareTestDatabase()

	t.Run("failures", func(t *testing.T) {
		assert := assert.New(t)

		items, err := CreateItems([]RequestPostItemsBody{}, "s9")
		assert.Error(err)
		assert.Empty(items)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)

		items, err := CreateItems([]RequestPostItemsBody{
			{Name: "test1", IsTrapItem: false, IsBook: false, Tags: []string{"test_tag", "test_tag2"},
				Description: "test_description", ImgURL: "https://example.com/"},
		}, "s9")
		assert.NoError(err)
		assert.NotEmpty(items)
		assert.NotEmpty(items[0].Ownership)

		items, err = CreateItems([]RequestPostItemsBody{
			{Name: "test2", IsTrapItem: false, IsBook: true, Tags: []string{"test_tag", "test_tag2"},
				Description: "test_description", ImgURL: "https://example.com/", Code: "9784088725093"},
			{Name: "test3", IsTrapItem: true, IsBook: false, Tags: []string{"test_tag", "test_tag2"},
				Description: "test_description", ImgURL: "https://example.com/", Count: 3},
			{Name: "test4", IsTrapItem: true, IsBook: true, Tags: []string{"test_tag", "test_tag2"},
				Description: "test_description", ImgURL: "https://example.com/", Code: "9784088725093", Count: 3},
		}, "s9")
		assert.NoError(err)
		assert.NotEmpty(items)

		assert.NotEmpty(items[0].Book)
		assert.Empty(items[0].Equipment)
		assert.NotEmpty(items[0].Ownership)

		assert.Empty(items[1].Book)
		assert.NotEmpty(items[1].Equipment)
		assert.Empty(items[1].Ownership)
	})
}

func TestGetItem(t *testing.T) {
	PrepareTestDatabase()
	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		res, err := GetItem(-1)
		assert.Error(err)
		assert.Empty(res)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		res, err := GetItem(1)
		assert.NoError(err)
		assert.NotEmpty(res)
		assert.Equal(res.Name, "item-id1")
	})
}

func TestPatchItem(t *testing.T) {
	PrepareTestDatabase()

	t.Run("failure", func(t *testing.T) {
		assert := assert.New(t)

		res, err := PatchItem(-1, RequestPostItemsBody{})
		assert.Error(err)
		assert.Empty(res)
	})

	t.Run("success-1", func(t *testing.T) {
		assert := assert.New(t)

		req := RequestPostItemsBody{
			Name: "testPatchItem", IsTrapItem: false, IsBook: false,
			Tags: []string{"tagTest"}, Description: "testPatchItem", ImgURL: "testURL",
		}
		res, err := PatchItem(1, req)

		assert.NoError(err)
		assert.NotEmpty(res)
		assert.Empty(res.Book)
		assert.Empty(res.Equipment)
		assert.Equal(res.Tag[0].Name, req.Tags[0])
		assert.Equal(res.Description, req.Description)
		assert.Equal(res.ImgURL, req.ImgURL)
		assert.NotEmpty(res.Ownership)
	})

	t.Run("success-2", func(t *testing.T) {
		assert := assert.New(t)

		req := RequestPostItemsBody{
			Name: "testPatchItem", IsTrapItem: true, IsBook: true, Count: 123456, Code: "9784088725093",
			Tags: []string{"tagTest"}, Description: "testPatchItem", ImgURL: "testURL",
		}
		res, err := PatchItem(4, req)

		assert.NoError(err)
		assert.NotEmpty(res)
		assert.Equal(res.Book.Code, req.Code)
		assert.Equal(res.Equipment.Count, req.Count)
		assert.Equal(res.Equipment.CountMax, req.Count)
		assert.Equal(res.Tag[0].Name, req.Tags[0])
		assert.Equal(res.Description, req.Description)
		assert.Equal(res.ImgURL, req.ImgURL)
	})
}

func TestDeleteItem(t *testing.T) {
	PrepareTestDatabase()

	t.Run("failure", func(t *testing.T) {
		assert := assert.New(t)
		err := DeleteItem(-1)
		assert.Error(err)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)
		err := DeleteItem(1)
		assert.NoError(err)

		err = DeleteItem(2)
		assert.NoError(err)
	})
}
