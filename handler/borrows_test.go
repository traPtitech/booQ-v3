package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/middleware"
	"github.com/traPtitech/booQ-v3/usecase"
	mock_usecase "github.com/traPtitech/booQ-v3/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestHandler_PostBorrow(t *testing.T) {
	dueDate := time.Date(2200, 7, 1, 23, 59, 59, 0, time.UTC)
	testCases := []struct {
		name         string
		itemID       string
		ownershipID  string
		userID       string
		requestBody  string
		setupMock    func(u *mock_usecase.MockBorrowingUseCase)
		expectedCode int
		expectedBody func() *openapi.BorrowRequest
	}{
		{
			name:        "success",
			itemID:      "1",
			ownershipID: "2",
			userID:      "user1",
			requestBody: `{"dueDate":"2200-07-01","propose":"for study","borrowInClubRoom":false}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					PostRequest("user1", 2, "for study", dueDate, false).
					Return(&domain.Transaction{
						ID:               100,
						UserID:           "user1",
						OwnershipID:      2,
						Purpose:          "for study",
						DueDate:          dueDate,
						BorrowInClubRoom: false,
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusCreated,
			expectedBody: func() *openapi.BorrowRequest {
				purpose := "for study"
				return &openapi.BorrowRequest{
					Propose:          &purpose,
					DueDate:          openapi_types.Date{Time: dueDate},
					BorrowInClubRoom: false,
				}
			},
		},
		{
			name:         "failure: invalid request body",
			itemID:       "1",
			ownershipID:  "2",
			userID:       "user1",
			requestBody:  `{`,
			setupMock:    func(u *mock_usecase.MockBorrowingUseCase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "failure: unauthorized",
			itemID:       "1",
			ownershipID:  "2",
			userID:       "",
			requestBody:  `{"dueDate":"2200-07-01","propose":"for study","borrowInClubRoom":false}`,
			setupMock:    func(u *mock_usecase.MockBorrowingUseCase) {},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:        "failure: usecase error",
			itemID:      "1",
			ownershipID: "2",
			userID:      "user1",
			requestBody: `{"dueDate":"2200-07-01","propose":"for study","borrowInClubRoom":false}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					PostRequest("user1", 2, "for study", dueDate, false).
					Return(nil, assert.AnError).
					Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "failure: ownership not found",
			itemID:      "1",
			ownershipID: "999",
			userID:      "user1",
			requestBody: `{"dueDate":"2200-07-01","propose":"for study","borrowInClubRoom":false}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					PostRequest("user1", 999, "for study", dueDate, false).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "failure: due date in the past",
			itemID:      "1",
			ownershipID: "2",
			userID:      "user1",
			requestBody: `{"dueDate":"2020-01-01","propose":"for study","borrowInClubRoom":false}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					PostRequest("user1", 2, "for study", time.Date(2020, 1, 1, 23, 59, 59, 0, time.UTC), false).
					Return(nil, usecase.ErrInvalidDueDate).
					Times(1)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBorrowingUseCase := mock_usecase.NewMockBorrowingUseCase(ctrl)
			tc.setupMock(mockBorrowingUseCase)

			h := NewHandler(nil, nil, nil, mockBorrowingUseCase)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/items/%s/owners/%s/borrowings", tc.itemID, tc.ownershipID), strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if tc.userID != "" {
				req = req.WithContext(middleware.WithUserID(req.Context(), tc.userID))
			}
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

func TestHandler_GetBorrowingById(t *testing.T) {
	dueDate := time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		name         string
		itemID       string
		ownershipID  string
		borrowingID  string
		userID       string
		setupMock    func(u *mock_usecase.MockBorrowingUseCase)
		expectedCode int
		expectedBody func() *openapi.Borrowing
	}{
		{
			name:        "success",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "user1",
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					GetRequest("user1", 2, 100).
					Return(&domain.Transaction{
						ID:               100,
						UserID:           "user1",
						OwnershipID:      2,
						Purpose:          "for study",
						DueDate:          dueDate,
						BorrowInClubRoom: false,
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() *openapi.Borrowing {
				purpose := "for study"
				return &openapi.Borrowing{
					Id:               100,
					Propose:          &purpose,
					DueDate:          openapi_types.Date{Time: dueDate},
					BorrowInClubRoom: false,
				}
			},
		},
		{
			name:        "failure: usecase error",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "user1",
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					GetRequest("user1", 2, 100).
					Return(nil, assert.AnError).
					Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "failure: unauthorized",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "",
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				// no calls expected
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:        "failure: ownership not found",
			itemID:      "1",
			ownershipID: "999",
			borrowingID: "100",
			userID:      "user1",
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					GetRequest("user1", 999, 100).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "failure: forbidden",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "user1",
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					GetRequest("user1", 2, 100).
					Return(nil, usecase.ErrForbidden).
					Times(1)
			},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBorrowingUseCase := mock_usecase.NewMockBorrowingUseCase(ctrl)
			tc.setupMock(mockBorrowingUseCase)

			h := NewHandler(nil, nil, nil, mockBorrowingUseCase)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/items/%s/owners/%s/borrowings/%s", tc.itemID, tc.ownershipID, tc.borrowingID), nil)
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

func TestHandler_PostBorrowReply(t *testing.T) {
	testCases := []struct {
		name         string
		itemID       string
		ownershipID  string
		borrowingID  string
		userID       string
		requestBody  string
		setupMock    func(u *mock_usecase.MockBorrowingUseCase)
		expectedCode int
		expectedBody func() *openapi.BorrowReply
	}{
		{
			name:        "success: approve",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "owner1",
			requestBody: `{"answer":true,"comment":"ok"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					ReplyRequest("owner1", 2, 100, true, "ok").
					Return(&domain.Transaction{
						Message: "ok",
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() *openapi.BorrowReply {
				return &openapi.BorrowReply{
					Answer:  true,
					Comment: "ok",
				}
			},
		},
		{
			name:        "success: reject",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "owner1",
			requestBody: `{"answer":false,"comment":"no"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					ReplyRequest("owner1", 2, 100, false, "no").
					Return(&domain.Transaction{
						Message: "no",
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() *openapi.BorrowReply {
				return &openapi.BorrowReply{
					Answer:  false,
					Comment: "no",
				}
			},
		},
		{
			name:        "failure: usecase error",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "owner1",
			requestBody: `{"answer":true,"comment":"ok"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					ReplyRequest("owner1", 2, 100, true, "ok").
					Return(nil, assert.AnError).
					Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "failure: unauthorized",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "",
			requestBody: `{"answer":true,"comment":"ok"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				// no calls expected
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:        "failure: ownership not found",
			itemID:      "1",
			ownershipID: "999",
			borrowingID: "100",
			userID:      "owner1",
			requestBody: `{"answer":true,"comment":"ok"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					ReplyRequest("owner1", 999, 100, true, "ok").
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "failure: forbidden",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "owner1",
			requestBody: `{"answer":true,"comment":"ok"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					ReplyRequest("owner1", 2, 100, true, "ok").
					Return(nil, usecase.ErrForbidden).
					Times(1)
			},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBorrowingUseCase := mock_usecase.NewMockBorrowingUseCase(ctrl)
			tc.setupMock(mockBorrowingUseCase)

			h := NewHandler(nil, nil, nil, mockBorrowingUseCase)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/items/%s/owners/%s/borrowings/%s/reply", tc.itemID, tc.ownershipID, tc.borrowingID), strings.NewReader(tc.requestBody))
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

func TestHandler_PostReturn(t *testing.T) {
	testCases := []struct {
		name         string
		itemID       string
		ownershipID  string
		borrowingID  string
		userID       string
		requestBody  string
		setupMock    func(u *mock_usecase.MockBorrowingUseCase)
		expectedCode int
	}{
		{
			name:        "success",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "user1",
			requestBody: `{"text":"returning"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					ReturnItem("user1", 2, 100, "returning").
					Return(nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "failure: usecase error",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "user1",
			requestBody: `{"text":"returning"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					ReturnItem("user1", 2, 100, "returning").
					Return(assert.AnError).
					Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "failure: unauthorized",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "",
			requestBody: `{"text":"returning"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				// no calls expected
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:        "failure: ownership not found",
			itemID:      "1",
			ownershipID: "999",
			borrowingID: "100",
			userID:      "user1",
			requestBody: `{"text":"returning"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					ReturnItem("user1", 999, 100, "returning").
					Return(domain.ErrNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "failure: forbidden",
			itemID:      "1",
			ownershipID: "2",
			borrowingID: "100",
			userID:      "user1",
			requestBody: `{"text":"returning"}`,
			setupMock: func(u *mock_usecase.MockBorrowingUseCase) {
				u.EXPECT().
					ReturnItem("user1", 2, 100, "returning").
					Return(usecase.ErrForbidden).
					Times(1)
			},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBorrowingUseCase := mock_usecase.NewMockBorrowingUseCase(ctrl)
			tc.setupMock(mockBorrowingUseCase)

			h := NewHandler(nil, nil, nil, mockBorrowingUseCase)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/items/%s/owners/%s/borrowings/%s/return", tc.itemID, tc.ownershipID, tc.borrowingID), strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if tc.userID != "" {
				req = req.WithContext(middleware.WithUserID(req.Context(), tc.userID))
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
