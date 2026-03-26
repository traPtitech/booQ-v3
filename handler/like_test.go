package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestHandler_AddLike(t *testing.T) {
	testCases := []struct {
		name         string
		itemID       string
		userID       string
		setupMock    func(u *mock_usecase.MockLikeUseCase)
		expectedCode int
	}{
		{
			name:   "success",
			itemID: "1",
			userID: "user1",
			setupMock: func(u *mock_usecase.MockLikeUseCase) {
				u.EXPECT().AddLike(1, "user1").Return(nil).Times(1)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:   "failure: unauthorized",
			itemID: "1",
			setupMock: func(u *mock_usecase.MockLikeUseCase) {
				// no calls expected
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:   "failure: already liked",
			itemID: "1",
			userID: "user1",
			setupMock: func(u *mock_usecase.MockLikeUseCase) {
				u.EXPECT().AddLike(1, "user1").Return(usecase.ErrAlreadyLiked).Times(1)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "failure: not found",
			itemID: "1",
			userID: "user1",
			setupMock: func(u *mock_usecase.MockLikeUseCase) {
				u.EXPECT().AddLike(1, "user1").Return(domain.ErrNotFound).Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLikeUseCase := mock_usecase.NewMockLikeUseCase(ctrl)
			tc.setupMock(mockLikeUseCase)

			h := NewHandlerWithTagLike(nil, nil, nil, nil, nil, mockLikeUseCase)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/items/%s/likes", tc.itemID), nil)
			if tc.userID != "" {
				req = req.WithContext(middleware.WithUserID(req.Context(), tc.userID))
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestHandler_RemoveLike(t *testing.T) {
	testCases := []struct {
		name         string
		itemID       string
		userID       string
		setupMock    func(u *mock_usecase.MockLikeUseCase)
		expectedCode int
	}{
		{
			name:   "success",
			itemID: "1",
			userID: "user1",
			setupMock: func(u *mock_usecase.MockLikeUseCase) {
				u.EXPECT().RemoveLike(1, "user1").Return(nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "failure: unauthorized",
			itemID: "1",
			setupMock: func(u *mock_usecase.MockLikeUseCase) {
				// no calls expected
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:   "failure: not liked",
			itemID: "1",
			userID: "user1",
			setupMock: func(u *mock_usecase.MockLikeUseCase) {
				u.EXPECT().RemoveLike(1, "user1").Return(usecase.ErrNotLiked).Times(1)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "failure: not found",
			itemID: "1",
			userID: "user1",
			setupMock: func(u *mock_usecase.MockLikeUseCase) {
				u.EXPECT().RemoveLike(1, "user1").Return(domain.ErrNotFound).Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLikeUseCase := mock_usecase.NewMockLikeUseCase(ctrl)
			tc.setupMock(mockLikeUseCase)

			h := NewHandlerWithTagLike(nil, nil, nil, nil, nil, mockLikeUseCase)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/items/%s/likes", tc.itemID), nil)
			if tc.userID != "" {
				req = req.WithContext(middleware.WithUserID(req.Context(), tc.userID))
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
