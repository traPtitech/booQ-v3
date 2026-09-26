package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/booQ-v3/internal/domain"
	"github.com/traPtitech/booQ-v3/internal/handler/openapi"
	"github.com/traPtitech/booQ-v3/internal/middleware"
	"github.com/traPtitech/booQ-v3/internal/usecase"
	mock_usecase "github.com/traPtitech/booQ-v3/internal/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestHandler_PostBorrowEquipment(t *testing.T) {
	due := time.Date(2200, 7, 1, 23, 59, 59, 0, time.UTC)
	const valid = `{"dueDate":"2200-07-01","propose":"purpose","count":3,"borrowInClubRoom":true}`
	for _, tc := range []struct {
		name, body   string
		unauth, skip bool
		count        int
		purpose      string
		err          error
		code         int
	}{
		{
			name:    "success",
			body:    valid,
			count:   3,
			purpose: "purpose",
			code:    http.StatusCreated,
		},
		{
			name:  "optional fields",
			body:  `{"dueDate":"2200-07-01"}`,
			count: 1,
			code:  http.StatusCreated,
		},
		{
			name: "malformed",
			body: `{`,
			skip: true,
			code: http.StatusBadRequest,
		},
		{
			name: "invalid date",
			body: `{"dueDate":"bad"}`,
			skip: true,
			code: http.StatusBadRequest,
		},
		{
			name: "missing date",
			body: `{}`,
			skip: true,
			code: http.StatusBadRequest,
		},
		{
			name:   "unauthorized",
			body:   valid,
			unauth: true,
			skip:   true,
			code:   http.StatusUnauthorized,
		},
		{
			name:    "missing item",
			body:    valid,
			count:   3,
			purpose: "purpose",
			err:     domain.ErrNotFound,
			code:    http.StatusNotFound,
		},
		{
			name:    "not equipment",
			body:    valid,
			count:   3,
			purpose: "purpose",
			err:     usecase.ErrItemNotEquipment,
			code:    http.StatusBadRequest,
		},
		{
			name:    "past due date",
			body:    valid,
			count:   3,
			purpose: "purpose",
			err:     usecase.ErrInvalidDueDate,
			code:    http.StatusBadRequest,
		},
		{
			name:    "internal error",
			body:    valid,
			count:   3,
			purpose: "purpose",
			err:     assert.AnError,
			code:    http.StatusInternalServerError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u := mock_usecase.NewMockBorrowingUseCase(gomock.NewController(t))

			if !tc.skip {
				var transaction *domain.EquipmentTransaction
				if tc.err == nil {
					transaction = domain.NewEquipmentTransaction("user", 7, tc.purpose, tc.count, due)
				}
				u.EXPECT().BorrowEquipment(7, "user", tc.purpose, tc.count, due).Return(transaction, tc.err)
			}

			req := httptest.NewRequest(http.MethodPost, "/items/7/borrowing/equipment", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if !tc.unauth {
				req = req.WithContext(middleware.WithUserID(req.Context(), "user"))
			}

			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)
			require.NoError(t, (&handler{bu: u}).PostBorrowEquipment(ctx, 7))

			assert.Equal(t, tc.code, rec.Code)
			if tc.code == 201 {
				if tc.name == "success" {
					assert.JSONEq(t, valid, rec.Body.String())
				} else {
					assert.JSONEq(t, `{"dueDate":"2200-07-01","propose":"","count":1,"borrowInClubRoom":false}`, rec.Body.String())
				}
			}
			if tc.code == 500 {
				assert.NotContains(t, rec.Body.String(), assert.AnError.Error())
			}
		})
	}
}

func TestHandler_PostBorrowEquipmentReturn(t *testing.T) {
	for _, tc := range []struct {
		name, body, message string
		unauth, skip        bool
		err                 error
		code                int
	}{
		{
			name:    "success",
			body:    `{"text":"returned"}`,
			message: "returned",
			code:    http.StatusCreated,
		},
		{
			name: "malformed",
			body: `{`,
			skip: true,
			code: http.StatusBadRequest,
		},
		{
			name:   "unauthorized",
			body:   `{"text":"returned"}`,
			unauth: true,
			skip:   true,
			code:   http.StatusUnauthorized,
		},
		{
			name: "empty text",
			body: `{"text":""}`,
			code: http.StatusCreated,
		},
		{
			name:    "different user",
			body:    `{"text":"returned"}`,
			message: "returned",
			err:     usecase.ErrForbidden,
			code:    http.StatusForbidden,
		},
		{
			name:    "missing transaction",
			body:    `{"text":"returned"}`,
			message: "returned",
			err:     domain.ErrNotFound,
			code:    http.StatusNotFound,
		},
		{
			name:    "invalid status",
			body:    `{"text":"returned"}`,
			message: "returned",
			err:     domain.ErrInvalidTransactionStatus,
			code:    http.StatusBadRequest,
		},
		{
			name:    "internal error",
			body:    `{"text":"returned"}`,
			message: "returned",
			err:     assert.AnError,
			code:    http.StatusInternalServerError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u := mock_usecase.NewMockBorrowingUseCase(gomock.NewController(t))

			if !tc.skip {
				u.EXPECT().ReturnEquipment(7, 4, "user", tc.message).Return(&domain.EquipmentTransaction{ReturnMessage: tc.message}, tc.err)
			}

			req := httptest.NewRequest(http.MethodPost, "/items/7/borrowing/equipment/4/return", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if !tc.unauth {
				req = req.WithContext(middleware.WithUserID(req.Context(), "user"))
			}

			rec := httptest.NewRecorder()
			e := echo.New()
			openapi.RegisterHandlers(e, &handler{bu: u})
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.code, rec.Code)
			if tc.code == 201 {
				assert.JSONEq(t, tc.body, rec.Body.String())
			}
			if tc.code == 500 {
				assert.NotContains(t, rec.Body.String(), assert.AnError.Error())
			}
		})
	}
}

func TestHandler_EquipmentReturnRoutes(t *testing.T) {
	for _, tc := range []struct {
		name, path string
		code       int
	}{
		{
			name: "invalid borrowing ID",
			path: "/items/7/borrowing/equipment/invalid/return",
			code: http.StatusBadRequest,
		},
		{
			name: "invalid item ID",
			path: "/items/invalid/borrowing/equipment/4/return",
			code: http.StatusBadRequest,
		},
		{
			name: "old route",
			path: "/items/7/borrowing/equipment/return",
			code: http.StatusNotFound,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u := mock_usecase.NewMockBorrowingUseCase(gomock.NewController(t))
			e := echo.New()
			openapi.RegisterHandlers(e, &handler{bu: u})

			req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(`{"text":"returned"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req = req.WithContext(middleware.WithUserID(req.Context(), "user"))
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.code, rec.Code)
		})
	}
}
