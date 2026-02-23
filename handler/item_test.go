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
	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/usecase"
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
			tc.setupMock(mockItemUseCase)

			h := NewHandler(mockItemUseCase)

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

func TestHandler_GetItems(t *testing.T) {
	testCases := []struct {
		name         string
		query        string
		setupMock    func(u *mock_usecase.MockItemUseCase)
		expectedCode int
		expectedBody []openapi.Item
	}{
		{
			name:  "success: no query",
			query: "",
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					SearchItems(domain.ItemSearchQuery{}).
					Return([]*domain.Item{
						{
							ID:          1,
							Name:        "Test Item 1",
							Description: "This is the first test item",
							ImgUrl:      "http://example.com/image1.png",
							BookDetail:  nil,
							EquipmentDetail: &domain.EquipmentDetail{
								Count:    1,
								CountMax: 5,
							},
							CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
							DeletedAt: nil,
						},
						{
							ID:          2,
							Name:        "Test Item 2",
							Description: "This is the second test item",
							ImgUrl:      "http://example.com/image2.png",
							BookDetail: &domain.BookDetail{
								ISBNCode: "0987654321",
							},
							EquipmentDetail: nil,
							CreatedAt:       time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
							UpdatedAt:       time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
							DeletedAt:       nil,
						},
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: []openapi.Item{
				{
					Id:          1,
					Name:        "Test Item 1",
					Description: "This is the first test item",
					ImgUrl:      "http://example.com/image1.png",
					IsBook:      false,
					IsTrapItem:  true,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nullable.NewNullNullable[time.Time](),
				},
				{
					Id:          2,
					Name:        "Test Item 2",
					Description: "This is the second test item",
					ImgUrl:      "http://example.com/image2.png",
					IsBook:      true,
					IsTrapItem:  false,
					CreatedAt:   time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nullable.NewNullNullable[time.Time](),
				},
			},
		},
		{
			name:  "failure: invalid query",
			query: "?limit=-1",
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					SearchItems(domain.ItemSearchQuery{Limit: -1}).
					Return(nil, usecase.ErrInvalidSearchQuery).
					Times(1)
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
			tc.setupMock(mockItemUseCase)

			h := NewHandler(mockItemUseCase)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/items%s", tc.query), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			body := strings.TrimSpace(rec.Body.String())

			if tc.expectedCode == http.StatusOK {
				expectedByte, err := json.Marshal(tc.expectedBody)
				assert.NoError(t, err)
				assert.Equal(t, string(expectedByte), body)
			}
		})
	}
}

func TestHandler_CreateItem(t *testing.T) {
	testCases := []struct {
		name         string
		requestBody  string
		setupMock    func(u *mock_usecase.MockItemUseCase)
		expectedCode int
		expectedBody *openapi.Item
	}{
		{
			name: "success",
			requestBody: `{
				"name": "New Item",
				"description": "This is a new item",
				"imgUrl": "http://example.com/new_image.png",
				"isBook": true,
				"isTrapItem": false
			}`,
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					CreateItem(&domain.Item{
						Name:        "New Item",
						Description: "This is a new item",
						ImgUrl:      "http://example.com/new_image.png",
						BookDetail: &domain.BookDetail{
							ISBNCode: "",
						},
						EquipmentDetail: nil,
					}).
					Return(&domain.Item{
						ID:          1,
						Name:        "New Item",
						Description: "This is a new item",
						ImgUrl:      "http://example.com/new_image.png",
						BookDetail: &domain.BookDetail{
							ISBNCode: "",
						},
						EquipmentDetail: nil,
						CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						DeletedAt:       nil,
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: &openapi.Item{
				Id:          1,
				Name:        "New Item",
				Description: "This is a new item",
				ImgUrl:      "http://example.com/new_image.png",
				IsBook:      true,
				IsTrapItem:  false,
				CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				DeletedAt:   nullable.NewNullNullable[time.Time](),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemUseCase := mock_usecase.NewMockItemUseCase(ctrl)
			tc.setupMock(mockItemUseCase)

			h := NewHandler(mockItemUseCase)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

func TestHandler_UpdateItem(t *testing.T) {
	testCases := []struct {
		name         string
		itemId       string
		requestBody  string
		setupMock    func(u *mock_usecase.MockItemUseCase)
		expectedCode int
		expectedBody *openapi.Item
	}{
		{
			name:   "success",
			itemId: "1",
			requestBody: `{
				"name": "Updated Item",
				"description": "This is an updated item",
				"imgUrl": "http://example.com/updated_image.png",
				"isBook": false,
				"isTrapItem": true
			}`,
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					UpdateItem(&domain.Item{
						ID:          1,
						Name:        "Updated Item",
						Description: "This is an updated item",
						ImgUrl:      "http://example.com/updated_image.png",
						BookDetail:  nil,
						EquipmentDetail: &domain.EquipmentDetail{
							Count:    0,
							CountMax: 0,
						},
					}).
					Return(&domain.Item{
						ID:          1,
						Name:        "Updated Item",
						Description: "This is an updated item",
						ImgUrl:      "http://example.com/updated_image.png",
						BookDetail:  nil,
						EquipmentDetail: &domain.EquipmentDetail{
							Count:    0,
							CountMax: 0,
						},
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
						DeletedAt: nil,
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: &openapi.Item{
				Id:          1,
				Name:        "Updated Item",
				Description: "This is an updated item",
				ImgUrl:      "http://example.com/updated_image.png",
				IsBook:      false,
				IsTrapItem:  true,
				CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				DeletedAt:   nullable.NewNullNullable[time.Time](),
			},
		},
		{
			name:   "failure: item not found",
			itemId: "2",
			requestBody: `{
				"name": "Updated Item",
				"description": "This is an updated item",
				"imgUrl": "http://example.com/updated_image.png",
				"isBook": false,
				"isTrapItem": true
			}`,
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					UpdateItem(&domain.Item{
						ID:          2,
						Name:        "Updated Item",
						Description: "This is an updated item",
						ImgUrl:      "http://example.com/updated_image.png",
						BookDetail:  nil,
						EquipmentDetail: &domain.EquipmentDetail{
							Count:    0,
							CountMax: 0,
						},
					}).
					Return(nil, domain.ErrItemNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: nil,
		},
		{
			name:   "failure: invalid item ID",
			itemId: "abc",
			requestBody: `{
				"name": "Updated Item",
				"description": "This is an updated item",
				"imgUrl": "http://example.com/updated_image.png",
				"isBook": false,
				"isTrapItem": true
			}`,
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				// No calls expected
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: nil,
		},
		{
			name:   "failure: cannot change item type",
			itemId: "1",
			requestBody: `{
				"name": "Updated Item",
				"description": "This is an updated item",
				"imgUrl": "http://example.com/updated_image.png",
				"isBook": false,
				"isTrapItem": false
			}`,
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					UpdateItem(&domain.Item{
						ID:              1,
						Name:            "Updated Item",
						Description:     "This is an updated item",
						ImgUrl:          "http://example.com/updated_image.png",
						BookDetail:      nil,
						EquipmentDetail: nil,
					}).
					Return(nil, fmt.Errorf("%w: cannot change whether item is equipment or not", usecase.ErrUpdateNotAllowed)).
					Times(1)
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
			tc.setupMock(mockItemUseCase)

			h := NewHandler(mockItemUseCase)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/items/%s", tc.itemId), strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

func TestHandler_DeleteItem(t *testing.T) {
	testCases := []struct {
		name         string
		itemId       string
		setupMock    func(u *mock_usecase.MockItemUseCase)
		expectedCode int
	}{
		{
			name:   "success",
			itemId: "1",
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					DeleteItem(1).
					Return(nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "failure: item not found",
			itemId: "2",
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				u.EXPECT().
					DeleteItem(2).
					Return(domain.ErrItemNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "failure: invalid item ID",
			itemId: "abc",
			setupMock: func(u *mock_usecase.MockItemUseCase) {
				// No calls expected
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemUseCase := mock_usecase.NewMockItemUseCase(ctrl)
			tc.setupMock(mockItemUseCase)

			h := NewHandler(mockItemUseCase)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/items/%s", tc.itemId), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
