package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/middleware"
	"github.com/traPtitech/booQ-v3/usecase"
	mock_usecase "github.com/traPtitech/booQ-v3/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestHandler_PostItemOwners(t *testing.T) {
	testCases := []struct {
		name         string
		itemID       string
		requestBody  string
		setupMock    func(u *mock_usecase.MockOwnershipUseCase)
		expectedCode int
		expectedBody func() *openapi.Ownership
	}{
		{
			name:   "success",
			itemID: "1",
			requestBody: `{
				"userId": "user1",
				"rentalable": true,
				"memo": "owner memo"
			}`,
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				u.EXPECT().
					CreateOwnership(&domain.Ownership{ItemID: 1, UserID: "user1", Rentable: true, Memo: "owner memo"}).
					Return(&domain.Ownership{ID: 10, ItemID: 1, UserID: "user1", Rentable: true, Memo: "owner memo"}, nil).
					Times(1)
			},
			expectedCode: http.StatusCreated,
			expectedBody: func() *openapi.Ownership {
				id := 10
				itemID := 1
				return &openapi.Ownership{Id: &id, ItemId: &itemID, UserId: "user1", Rentalable: true, Memo: "owner memo"}
			},
		},
		{
			name:        "failure: invalid request body",
			itemID:      "1",
			requestBody: `{`,
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				// no calls expected
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "failure: usecase error",
			itemID: "1",
			requestBody: `{
				"userId": "user1",
				"rentalable": true,
				"memo": "owner memo"
			}`,
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				u.EXPECT().
					CreateOwnership(&domain.Ownership{ItemID: 1, UserID: "user1", Rentable: true, Memo: "owner memo"}).
					Return(nil, assert.AnError).
					Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOwnershipUseCase := mock_usecase.NewMockOwnershipUseCase(ctrl)
			tc.setupMock(mockOwnershipUseCase)

			h := NewHandler(nil, nil, mockOwnershipUseCase)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/items/%s/owners", tc.itemID), strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedCode == http.StatusCreated {
				expectedBody, err := json.Marshal(tc.expectedBody())
				assert.NoError(t, err)
				assert.Equal(t, string(expectedBody), strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}

func TestHandler_EditItemOwners(t *testing.T) {
	testCases := []struct {
		name         string
		itemID       string
		ownershipID  string
		requestBody  string
		userID       string
		setupMock    func(u *mock_usecase.MockOwnershipUseCase)
		expectedCode int
		expectedBody func() *openapi.Ownership
	}{
		{
			name:        "success",
			itemID:      "1",
			ownershipID: "10",
			userID:      "owner",
			requestBody: `{"userId":"owner","rentalable":false,"memo":"updated memo"}`,
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				u.EXPECT().
					UpdateOwnership(&domain.Ownership{ID: 10, ItemID: 1, UserID: "owner", Rentable: false, Memo: "updated memo"}, "owner").
					Return(&domain.Ownership{ID: 10, ItemID: 1, UserID: "owner", Rentable: false, Memo: "updated memo"}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() *openapi.Ownership {
				id := 10
				itemID := 1
				return &openapi.Ownership{Id: &id, ItemId: &itemID, UserId: "owner", Rentalable: false, Memo: "updated memo"}
			},
		},
		{
			name:        "failure: invalid request body",
			itemID:      "1",
			ownershipID: "10",
			userID:      "owner",
			requestBody: `{`,
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				// no calls expected
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "failure: unauthorized",
			itemID:      "1",
			ownershipID: "10",
			requestBody: `{"userId":"owner","rentalable":false,"memo":"updated memo"}`,
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				// no calls expected
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:        "failure: forbidden",
			itemID:      "1",
			ownershipID: "10",
			userID:      "another-user",
			requestBody: `{"userId":"owner","rentalable":false,"memo":"updated memo"}`,
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				u.EXPECT().
					UpdateOwnership(&domain.Ownership{ID: 10, ItemID: 1, UserID: "owner", Rentable: false, Memo: "updated memo"}, "another-user").
					Return(nil, usecase.ErrForbidden).
					Times(1)
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name:        "failure: not found",
			itemID:      "1",
			ownershipID: "10",
			userID:      "owner",
			requestBody: `{"userId":"owner","rentalable":false,"memo":"updated memo"}`,
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				u.EXPECT().
					UpdateOwnership(&domain.Ownership{ID: 10, ItemID: 1, UserID: "owner", Rentable: false, Memo: "updated memo"}, "owner").
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOwnershipUseCase := mock_usecase.NewMockOwnershipUseCase(ctrl)
			tc.setupMock(mockOwnershipUseCase)

			h := NewHandler(nil, nil, mockOwnershipUseCase)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/items/%s/owners/%s", tc.itemID, tc.ownershipID), strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if tc.userID != "" {
				req = req.WithContext(middleware.WithUserID(req.Context(), tc.userID))
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedCode == http.StatusOK {
				expectedBody, err := json.Marshal(tc.expectedBody())
				assert.NoError(t, err)
				assert.Equal(t, string(expectedBody), strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}

func TestHandler_DeleteItemOwners(t *testing.T) {
	testCases := []struct {
		name         string
		itemID       string
		ownershipID  string
		userID       string
		setupMock    func(u *mock_usecase.MockOwnershipUseCase)
		expectedCode int
	}{
		{
			name:        "success",
			itemID:      "1",
			ownershipID: "10",
			userID:      "owner",
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				u.EXPECT().DeleteOwnership(10, 1, "owner").Return(nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "failure: unauthorized",
			itemID:      "1",
			ownershipID: "10",
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				// no calls expected
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:        "failure: forbidden",
			itemID:      "1",
			ownershipID: "10",
			userID:      "another-user",
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				u.EXPECT().DeleteOwnership(10, 1, "another-user").Return(usecase.ErrForbidden).Times(1)
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name:        "failure: not found",
			itemID:      "1",
			ownershipID: "10",
			userID:      "owner",
			setupMock: func(u *mock_usecase.MockOwnershipUseCase) {
				u.EXPECT().DeleteOwnership(10, 1, "owner").Return(domain.ErrNotFound).Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOwnershipUseCase := mock_usecase.NewMockOwnershipUseCase(ctrl)
			tc.setupMock(mockOwnershipUseCase)

			h := NewHandler(nil, nil, mockOwnershipUseCase)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/items/%s/owners/%s", tc.itemID, tc.ownershipID), nil)
			if tc.userID != "" {
				req = req.WithContext(middleware.WithUserID(req.Context(), tc.userID))
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
