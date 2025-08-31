package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/model"
)

func TestPostOwners(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider(TEST_USER))

	testCase := []struct {
		name         string
		itemID       string
		payload      model.OwnershipPayload
		expectedCode int
	}{
		{
			name:         "failure: invalid item id",
			itemID:       "aaa",
			payload:      model.OwnershipPayload{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "failure: item not found",
			itemID:       "999",
			payload:      model.OwnershipPayload{},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "success",
			itemID: "1",
			payload: model.OwnershipPayload{
				UserID:     "new_user",
				Rentalable: true,
				Memo:       "memo 1",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(echo.POST, fmt.Sprintf("/api/items/%s/owners", tc.itemID), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusOK {
				itemID, _ := strconv.Atoi(tc.itemID)

				res := model.Ownership{}
				_ = json.NewDecoder(rec.Body).Decode(&res)
				assert.NotZero(t, res.ID)
				assert.Equal(t, itemID, res.ItemID)
				assert.Equal(t, tc.payload.UserID, res.UserID)
				assert.Equal(t, tc.payload.Rentalable, res.Rentalable)
				assert.Equal(t, tc.payload.Memo, res.Memo)
			}
		})
	}
}

func TestPatchOwners(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider(TEST_USER))

	testCase := []struct {
		name         string
		itemID       string
		ownershipID  string
		payload      model.OwnershipPayload
		expectedCode int
	}{
		{
			name:         "failure: bad item id",
			itemID:       "aaa",
			ownershipID:  "1",
			payload:      model.OwnershipPayload{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "failure: bad ownership id",
			itemID:       "1",
			ownershipID:  "bbb",
			payload:      model.OwnershipPayload{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "failure: unauthorized update",
			itemID:       "1",
			ownershipID:  "1",
			payload:      model.OwnershipPayload{UserID: "cp20", Rentalable: false, Memo: "hack"},
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "success",
			itemID:       "1",
			ownershipID:  "1",
			payload:      model.OwnershipPayload{UserID: "s9", Rentalable: false, Memo: "updated"},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(echo.PATCH, fmt.Sprintf("/api/items/%s/owners/%s", tc.itemID, tc.ownershipID), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusOK {
				res := model.Ownership{}
				_ = json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(t, tc.payload.UserID, res.UserID)
				assert.Equal(t, tc.payload.Rentalable, res.Rentalable)
				assert.Equal(t, tc.payload.Memo, res.Memo)
			}
		})
	}
}

func TestDeleteOwners(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider(TEST_USER))

	testCase := []struct {
		name         string
		itemID       string
		ownershipID  string
		expectedCode int
	}{
		{
			name:         "failure: bad ownership id",
			itemID:       "1",
			ownershipID:  "aaa",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "failure: unauthorized delete",
			itemID:       "1",
			ownershipID:  "2", // ownership 2 not belongs to TEST_USER
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "success",
			itemID:       "1",
			ownershipID:  "1",
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, fmt.Sprintf("/api/items/%s/owners/%s", tc.itemID, tc.ownershipID), nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusNoContent {
				ownID, _ := strconv.Atoi(tc.ownershipID)
				own, err := model.GetOwnership(ownID)
				assert.ErrorIs(t, err, model.ErrNotFound)
				assert.Empty(t, own)
			}
		})
	}
}
