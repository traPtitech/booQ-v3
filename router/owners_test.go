package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/model"
)

/*
modelの方のテストをrouterに移動する

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

*/

/*
type PostOwnershipBody struct {
	Rentalable bool   `json:"rentalable"`
	Memo       string `json:"memo"`
}
*/

func TestPostOwners(t *testing.T) {
	model.PrepareTestDatabase()
	e := echo.New()
	SetupRouting(e, CreateUserProvider("s9"))

	t.Run("failure-invaild-id", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.PostOwnershipBody{
			Rentalable: true, Memo: "test"}

		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest("POST", "/api/items/-1/owners", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("failure-setting-owner-to-equipment", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.PostOwnershipBody{
			Rentalable: true, Memo: "test"}

		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest("POST", "/api/items/3/owners", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.PostOwnershipBody{
			Rentalable: true, Memo: "test"}

		reqBody, _ := json.Marshal(reqStruct)
		fmt.Print(string(reqBody))
		req := httptest.NewRequest("POST", "/api/items/1/owners", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(http.StatusOK, rec.Code)

		fmt.Println(rec.Body.String())
		ownership := model.Ownership{}
		_ = json.NewDecoder(rec.Body).Decode(&ownership)

		assert.Equal(http.StatusOK, rec.Code)
		assert.Equal(ownership.GormModel.ID, 5)
		assert.Equal(ownership.ItemID, 1)
		assert.Equal(ownership.UserID, "s9")
		assert.Equal(ownership.Rentalable, true)
		assert.Equal(ownership.Memo, "test")
	})
}
func TestPatchOwners(t *testing.T) {
	model.PrepareTestDatabase()
	e := echo.New()
	SetupRouting(e, CreateUserProvider("s9"))

	t.Run("failure-invaild-item-id", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.PostOwnershipBody{
			Rentalable: true, Memo: "test",
		}
		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest("PATCH", "/api/items/-1/owners/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("failure-invaild-ownership-id", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.PostOwnershipBody{
			Rentalable: true, Memo: "test",
		}
		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest("PATCH", "/api/items/1/owners/-1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)

		reqStruct := model.PostOwnershipBody{
			Rentalable: true, Memo: "test",
		}
		reqBody, _ := json.Marshal(reqStruct)
		req := httptest.NewRequest("PATCH", "/api/items/1/owners/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		fmt.Println(rec.Body.String())

		ownership := model.Ownership{}
		_ = json.NewDecoder(rec.Body).Decode(&ownership)

		assert.Equal(http.StatusOK, rec.Code)
		assert.Equal(ownership.ItemID, 1)
		assert.Equal(ownership.UserID, "s9")
		assert.Equal(ownership.Rentalable, true)
		assert.Equal(ownership.Memo, "test")
	})
}
func TestDeleteOwners(t *testing.T) {
	model.PrepareTestDatabase()
	e := echo.New()
	SetupRouting(e, CreateUserProvider("s9"))
}
