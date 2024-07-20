package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/model"
)

func TestGetItems(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider("s9"))

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		req := httptest.NewRequest(echo.GET, "/api/items?limit=aaa", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusBadRequest, rec.Code)
	})
	t.Run("success-1", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		req := httptest.NewRequest(echo.GET, "/api/items?search=item-id4", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		items := []model.Item{}
		_ = json.NewDecoder(rec.Body).Decode(&items)

		assert.Equal("item-id4 book equipment", items[0].Name)
		assert.Equal("aaa", items[0].Description)
		assert.Equal("url", items[0].ImgURL)
		assert.Equal(90, items[0].Equipment.Count)
		assert.Equal(100, items[0].Equipment.CountMax)
		assert.Equal("9784088725093", items[0].Book.Code)
		assert.Equal("tag3", items[0].Tag[0].Name)
		// 他の情報は/itemsではなく/items/{id}で取得させる
	})
}

func TestPostItems(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider("s9"))

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := []model.RequestPostItemsBody{
			{Name: "item-id5", IsTrapItem: false, IsBook: false, Tags: []string{"tagtest"},
				Description: "bbb", ImgURL: "url"},
			{Name: "item-id6", IsTrapItem: true, IsBook: true, Tags: []string{"tagtest"},
				Description: "bbb", ImgURL: "url", Count: 100, Code: "9784041026403"},
		}

		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest(echo.POST, "/api/items", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		items := []model.Item{}
		_ = json.NewDecoder(rec.Body).Decode(&items)

		assert.Equal(reqStruct[0].Name, items[0].Name)
		assert.Empty(items[0].Book)
		assert.Empty(items[0].Equipment)
		assert.Equal(reqStruct[0].Description, items[0].Description)
		assert.Equal(reqStruct[0].ImgURL, items[0].ImgURL)
		assert.Equal("s9", items[0].Ownership[0].UserID)

		assert.Equal(reqStruct[1].Name, items[1].Name)
		assert.Equal(reqStruct[1].Code, items[1].Book.Code)
		assert.Equal(reqStruct[1].Count, items[1].Equipment.Count)
		assert.Equal(reqStruct[1].Count, items[1].Equipment.CountMax)
		assert.Equal(reqStruct[1].Description, items[1].Description)
		assert.Equal(reqStruct[1].ImgURL, items[1].ImgURL)
		assert.Empty(items[1].Ownership)
	})
}

func TestGetItem(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider("s9"))

	t.Run("failure-1", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		req := httptest.NewRequest(echo.GET, "/api/items/-1", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(http.StatusInternalServerError, rec.Code)
	})

	t.Run("failure-2", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		req := httptest.NewRequest(echo.GET, "/api/items/aaa", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("success-1", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		req := httptest.NewRequest(echo.GET, "/api/items/1", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		item := model.Item{}
		_ = json.NewDecoder(rec.Body).Decode(&item)

		assert.Equal("item-id1", item.Name)
		assert.Equal("aaa", item.Description)
		assert.Equal("url", item.ImgURL)
		assert.Empty(item.Equipment)
		assert.Empty(item.Book)
		assert.Equal("tag1", item.Tag[0].Name)
		assert.Empty(item.TransactionEquipment)
		assert.Equal("s9", item.Comment[0].UserID)
		assert.Equal("comment", item.Comment[0].Comment)
		assert.Equal("cp20", item.Like[0].UserID)
		assert.Equal("s9", item.Ownership[0].UserID)
		assert.Equal(true, item.Ownership[0].Rentalable)
		assert.Equal("memo1", item.Ownership[0].Memo)
		assert.Equal("ryoha", item.Ownership[0].Transaction[2].UserID)
		assert.Equal("かりたいから", item.Ownership[0].Transaction[2].Purpose)
		assert.Equal("いいよ", item.Ownership[0].Transaction[2].Message)
		assert.Equal("ありがとう", item.Ownership[0].Transaction[2].ReturnMessage)
	})

	t.Run("success-2", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		req := httptest.NewRequest(echo.GET, "/api/items/4", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		item := model.Item{}
		_ = json.NewDecoder(rec.Body).Decode(&item)

		assert.Equal("item-id4 book equipment", item.Name)
		assert.Equal("aaa", item.Description)
		assert.Equal("url", item.ImgURL)
		assert.Equal(90, item.Equipment.Count)
		assert.Equal(100, item.Equipment.CountMax)
		assert.Equal("9784088725093", item.Book.Code)
		assert.Equal("tag3", item.Tag[0].Name)
		assert.Equal("ryoha", item.TransactionEquipment[0].UserID)
		assert.Equal("かりたいから", item.TransactionEquipment[0].Purpose)
		assert.Equal("かえしました", item.TransactionEquipment[0].ReturnMessage)
	})
}

func TestPatchItem(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider("s9"))

	t.Run("failure-1", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.RequestPostItemsBody{
			Name: "item-id5", IsTrapItem: false, IsBook: false, Tags: []string{"tagtest"},
			Description: "bbb", ImgURL: "url"}

		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest(echo.PATCH, "/api/items/4", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusInternalServerError, rec.Code)
	})

	t.Run("failure-2", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.RequestPostItemsBody{
			Name: "item-id6", IsTrapItem: true, IsBook: true, Tags: []string{"tagtest"},
			Description: "bbb", ImgURL: "url", Count: 100, Code: "9784041026403"}

		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest(echo.PATCH, "/api/items/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusInternalServerError, rec.Code)
	})

	t.Run("success-1", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.RequestPostItemsBody{
			Name: "item-id5", IsTrapItem: false, IsBook: false, Tags: []string{"tagtest"},
			Description: "bbb", ImgURL: "url"}

		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest(echo.PATCH, "/api/items/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		item := model.Item{}
		_ = json.NewDecoder(rec.Body).Decode(&item)

		assert.Equal(reqStruct.Name, item.Name)
		assert.Empty(item.Book)
		assert.Empty(item.Equipment)
		assert.Equal(reqStruct.Description, item.Description)
		assert.Equal(reqStruct.ImgURL, item.ImgURL)
		assert.Equal("s9", item.Ownership[0].UserID)
	})

	t.Run("success-2", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.RequestPostItemsBody{
			Name: "item-id6", IsTrapItem: true, IsBook: true, Tags: []string{"tagtest"},
			Description: "bbb", ImgURL: "url", Count: 111, Code: "9784123457689"}

		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest(echo.PATCH, "/api/items/4", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		item := model.Item{}
		_ = json.NewDecoder(rec.Body).Decode(&item)

		assert.Equal(reqStruct.Name, item.Name)
		assert.Equal("9784123457689", item.Book.Code)
		assert.Equal(111, item.Equipment.Count)
		assert.Equal(111, item.Equipment.CountMax)
		assert.Equal(reqStruct.Description, item.Description)
		assert.Equal(reqStruct.ImgURL, item.ImgURL)
		assert.Empty(item.Ownership)
	})
}

func TestDeleteItem(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider("s9"))

	t.Run("failure-1", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.DELETE, "/api/items/-1", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(http.StatusInternalServerError, rec.Code)
	})

	t.Run("failure-2", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.DELETE, "/api/items/aaa", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("success-1", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.DELETE, "/api/items/1", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		item, err := model.GetItem(1)
		assert.Error(err)
		assert.Empty(item)
	})

	t.Run("success-2", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.DELETE, "/api/items/4", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		item, err := model.GetItem(4)
		assert.Error(err)
		assert.Empty(item)
	})
}
