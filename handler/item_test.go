package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	mock_usecase "github.com/traPtitech/booQ-v3/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestHandler_GetItem(t *testing.T) {
	testCases := []struct {
		name         string
		itemId       string
		setupMock    func(u *mock_usecase.MockItemUseCase)
		expectedCode int
		expectedBody *openapi.Item
	}{
		{
			name:   "success",
			itemId: "1",
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					GetItemByID(1).
					Return(&domain.Item{
						ID:          1,
						Name:        "Test Item",
						Description: "This is a test item",
						ImgUrl:      "http://example.com/image.png",
						BookDetail: &domain.BookDetail{
							ISBNCode: "1234567890",
						},
						EquipmentDetail: nil,
						CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:       time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
						DeletedAt:       nil,
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: &openapi.Item{
				Id:          1,
				Name:        "Test Item",
				Description: "This is a test item",
				ImgUrl:      "http://example.com/image.png",
				IsBook:      true,
				IsTrapItem:  false,
				CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				DeletedAt:   nullable.NewNullNullable[time.Time](),
			},
		},
		{
			name:   "failure: item not found",
			itemId: "2",
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					GetItemByID(2).
					Return(nil, domain.ErrItemNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: nil,
		},
		{
			name:   "failure: invalid item ID",
			itemId: "abc",
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				// No calls expected
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemUseCase := mock_usecase.NewMockItemUseCase(ctrl)
			mockFileUseCase := mock_usecase.NewMockFileUseCase(ctrl)
			tc.setupMock(mockItemUseCase)

			h := NewHandler(mockItemUseCase, mockFileUseCase)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/items/%s", tc.itemId), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			body := strings.TrimSpace(rec.Body.String())

			if tc.expectedCode == http.StatusOK {
				expectedByte, err := tc.expectedBody.MarshalJSON()
				assert.NoError(t, err)
				assert.Equal(t, string(expectedByte), body)
			}
		})
	}
}
