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

func TestPostOwners(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider(TEST_USER))

	t.Run("failure: invalid item id", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.POST, "/api/items/aaa/owners", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)

		payload := model.OwnershipPayload{
			UserID:     "new_user",
			Rentalable: true,
			Memo:       "memo 1",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(echo.POST, "/api/items/1/owners", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		res := model.Ownership{}
		_ = json.NewDecoder(rec.Body).Decode(&res)
		assert.NotZero(res.ID)
		assert.Equal(1, res.ItemID)
		assert.Equal(payload.UserID, res.UserID)
		assert.Equal(payload.Rentalable, res.Rentalable)
		assert.Equal(payload.Memo, res.Memo)
	})
}

func TestPatchOwners(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider(TEST_USER))

	t.Run("failure: bad item id", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.PATCH, "/api/items/aaa/owners/1", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("failure: bad ownership id", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.PATCH, "/api/items/1/owners/bbb", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("failure: unauthorized update", func(t *testing.T) {
		assert := assert.New(t)
		payload := model.OwnershipPayload{UserID: "cp20", Rentalable: false, Memo: "hack"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(echo.PATCH, "/api/items/1/owners/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		// model.UpdateOwnership returns error -> router maps to 500
		assert.Equal(http.StatusInternalServerError, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)
		payload := model.OwnershipPayload{UserID: "s9", Rentalable: false, Memo: "updated"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(echo.PATCH, "/api/items/1/owners/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)

		res := model.Ownership{}
		_ = json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal("s9", res.UserID)
		assert.Equal(false, res.Rentalable)
		assert.Equal("updated", res.Memo)
	})
}

func TestDeleteOwners(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider(TEST_USER))

	t.Run("failure: bad ownership id", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.DELETE, "/api/items/1/owners/aaa", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("failure: unauthorized delete", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.DELETE, "/api/items/1/owners/2", nil) // ownership 2 belongs to cp20
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(http.StatusInternalServerError, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest(echo.DELETE, "/api/items/1/owners/1", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusNoContent, rec.Code)

		// Ensure deleted
		own, err := model.GetOwnership(1)
		assert.Error(err)
		assert.Empty(own)
	})
}
