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
	openapi_types "github.com/oapi-codegen/runtime/types"
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
		setupMock    func(iu *mock_usecase.MockItemUseCase)
		expectedCode int
		expectedBody func() *openapi.ItemDetail
	}{
		{
			name:   "success",
			itemId: "1",
			setupMock: func(iu *mock_usecase.MockItemUseCase) {
				iu.EXPECT().
					GetItemDetailByID(1).
					Return(&domain.ItemDetail{
						Item: &domain.Item{
							ID:              1,
							Name:            "Test Item",
							Description:     "This is a test item",
							ImgUrl:          "http://example.com/image.png",
							EquipmentDetail: nil,
							CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:       time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
							DeletedAt:       nil,
						},
						Tags: []*domain.Tag{{Name: "tag-1"}},
						Likes: []*domain.Like{
							{ItemID: 1, UserID: "user1"},
							{ItemID: 1, UserID: "user2"},
						},
						Ownerships: []*domain.OwnershipDetail{
							{
								Ownership: &domain.Ownership{
									ID:       20,
									ItemID:   1,
									UserID:   "owner1",
									Rentable: true,
								},
								Transactions: []*domain.Transaction{
									{
										ID:          10,
										UserID:      "borrower1",
										OwnershipID: 20,
										Status:      domain.BorrowingStatusRequested,
										Purpose:     "read",
										DueDate:     time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
										CreatedAt:   time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
										UpdatedAt:   time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
									},
								},
							},
							{
								Ownership: &domain.Ownership{
									ID:       21,
									ItemID:   1,
									UserID:   "owner2",
									Rentable: true,
								},
								Transactions: []*domain.Transaction{
									{
										ID:          11,
										UserID:      "borrower2",
										OwnershipID: 21,
										Status:      domain.BorrowingStatusBorrowed,
										Purpose:     "research",
										DueDate:     time.Date(2025, 1, 11, 0, 0, 0, 0, time.UTC),
										CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
										UpdatedAt:   time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
									},
								},
							},
						},
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() *openapi.ItemDetail {
				transactionCreatedAt := time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)
				transactionID := 10
				ownershipID := 20
				purpose := "read"
				status := 0
				transactionUpdatedAt := time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC)
				userID := "borrower1"
				transactionCreatedAt2 := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
				transactionID2 := 11
				ownershipID2 := 21
				purpose2 := "research"
				status2 := 1
				transactionUpdatedAt2 := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
				userID2 := "borrower2"

				res := &openapi.ItemDetail{
					Id:          1,
					Name:        "Test Item",
					Description: "This is a test item",
					ImgUrl:      "http://example.com/image.png",
					IsBook:      false,
					IsTrapItem:  false,
					LikesByUsers: &[]string{
						"user1",
						"user2",
					},
					Tags: &[]openapi.Tag{
						{Name: "tag-1"},
					},
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					DeletedAt: nullable.NewNullNullable[time.Time](),
				}
				_ = res.FromItemDetail1(openapi.ItemDetail1{
					Transactions: []openapi.Transaction{
						{
							CheckoutDate:  openapi_types.Date{},
							CreatedAt:     &transactionCreatedAt,
							DeletedAt:     nullable.NewNullNullable[time.Time](),
							DueDate:       openapi_types.Date{Time: time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)},
							Id:            &transactionID,
							Message:       "",
							OwnershipId:   &ownershipID,
							Purpose:       &purpose,
							ReturnMessage: "",
							ReturnDate:    openapi_types.Date{},
							Status:        &status,
							UpdatedAt:     &transactionUpdatedAt,
							UserId:        &userID,
						},
						{
							CheckoutDate:  openapi_types.Date{},
							CreatedAt:     &transactionCreatedAt2,
							DeletedAt:     nullable.NewNullNullable[time.Time](),
							DueDate:       openapi_types.Date{Time: time.Date(2025, 1, 11, 0, 0, 0, 0, time.UTC)},
							Id:            &transactionID2,
							Message:       "",
							OwnershipId:   &ownershipID2,
							Purpose:       &purpose2,
							ReturnMessage: "",
							ReturnDate:    openapi_types.Date{},
							Status:        &status2,
							UpdatedAt:     &transactionUpdatedAt2,
							UserId:        &userID2,
						},
					},
				})
				return res
			},
		},
		{
			name:   "success: item is book",
			itemId: "1",
			setupMock: func(iu *mock_usecase.MockItemUseCase) {
				iu.EXPECT().
					GetItemDetailByID(1).
					Return(&domain.ItemDetail{
						Item: &domain.Item{
							ID:          1,
							Name:        "Test Book",
							Description: "This is a test book",
							ImgUrl:      "http://example.com/book.png",
							BookDetail: &domain.BookDetail{
								ISBNCode: "1234567890",
							},
							EquipmentDetail: nil,
							CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:       time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
							DeletedAt:       nil,
						},
						Tags:       []*domain.Tag{},
						Likes:      []*domain.Like{},
						Ownerships: []*domain.OwnershipDetail{},
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() *openapi.ItemDetail {
				code := "1234567890"

				res := &openapi.ItemDetail{
					Id:          1,
					Name:        "Test Book",
					Description: "This is a test book",
					ImgUrl:      "http://example.com/book.png",
					IsBook:      true,
					IsTrapItem:  false,
					Code:        &code,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nullable.NewNullNullable[time.Time](),
				}
				_ = res.FromItemDetail1(openapi.ItemDetail1{Transactions: []openapi.Transaction{}})
				return res
			},
		},
		{
			name:   "success: item is equipment",
			itemId: "1",
			setupMock: func(iu *mock_usecase.MockItemUseCase) {
				iu.EXPECT().
					GetItemDetailByID(1).
					Return(&domain.ItemDetail{
						Item: &domain.Item{
							ID:          1,
							Name:        "Test Equipment",
							Description: "This is a test equipment",
							ImgUrl:      "http://example.com/equipment.png",
							BookDetail:  nil,
							EquipmentDetail: &domain.EquipmentDetail{
								Count:    2,
								CountMax: 5,
							},
							CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
							DeletedAt: nil,
						},
						Tags:  []*domain.Tag{{Name: "equipment"}},
						Likes: []*domain.Like{},
						Ownerships: []*domain.OwnershipDetail{
							{
								Ownership: &domain.Ownership{
									ID:       21,
									ItemID:   1,
									UserID:   "owner2",
									Rentable: true,
								},
								Transactions: []*domain.Transaction{
									{
										ID:          11,
										UserID:      "borrower2",
										OwnershipID: 21,
										Status:      domain.BorrowingStatusBorrowed,
										Purpose:     "use",
										DueDate:     time.Date(2025, 1, 11, 0, 0, 0, 0, time.UTC),
										CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
										UpdatedAt:   time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
									},
								},
							},
						},
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() *openapi.ItemDetail {
				count := 2
				transactionCreatedAt := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
				transactionID := 11
				itemID := 1
				purpose := "use"
				status := 1
				transactionUpdatedAt := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
				userID := "borrower2"

				item := &openapi.ItemDetail{
					Id:          1,
					Name:        "Test Equipment",
					Description: "This is a test equipment",
					ImgUrl:      "http://example.com/equipment.png",
					Count:       &count,
					IsBook:      false,
					IsTrapItem:  true,
					Tags: &[]openapi.Tag{
						{Name: "equipment"},
					},
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					DeletedAt: nullable.NewNullNullable[time.Time](),
				}
				_ = item.FromItemDetail0(openapi.ItemDetail0{
					TransactionsEquipment: []openapi.TransactionEquipment{
						{
							CreatedAt:  &transactionCreatedAt,
							DeletedAt:  nullable.NewNullNullable[time.Time](),
							DueDate:    openapi_types.Date{Time: time.Date(2025, 1, 11, 0, 0, 0, 0, time.UTC)},
							Id:         &transactionID,
							ItemId:     &itemID,
							Purpose:    &purpose,
							ReturnDate: openapi_types.Date{},
							Status:     &status,
							UpdatedAt:  &transactionUpdatedAt,
							UserId:     &userID,
						},
					},
				})
				return item
			},
		},
		{
			name:   "success: item is book and equipment",
			itemId: "1",
			setupMock: func(iu *mock_usecase.MockItemUseCase) {
				iu.EXPECT().
					GetItemDetailByID(1).
					Return(&domain.ItemDetail{
						Item: &domain.Item{
							ID:          1,
							Name:        "Test Book Equipment",
							Description: "This is a test book equipment",
							ImgUrl:      "http://example.com/book_equipment.png",
							BookDetail: &domain.BookDetail{
								ISBNCode: "1234567890",
							},
							EquipmentDetail: &domain.EquipmentDetail{
								Count:    2,
								CountMax: 5,
							},
							CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
							DeletedAt: nil,
						},
						Tags:       []*domain.Tag{},
						Likes:      []*domain.Like{},
						Ownerships: []*domain.OwnershipDetail{},
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() *openapi.ItemDetail {
				code := "1234567890"
				count := 2
				item := &openapi.ItemDetail{
					Id:          1,
					Name:        "Test Book Equipment",
					Description: "This is a test book equipment",
					ImgUrl:      "http://example.com/book_equipment.png",
					Count:       &count,
					IsBook:      true,
					IsTrapItem:  true,
					Code:        &code,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nullable.NewNullNullable[time.Time](),
				}
				_ = item.FromItemDetail0(openapi.ItemDetail0{TransactionsEquipment: []openapi.TransactionEquipment{}})
				return item
			},
		},
		{
			name:   "failure: item not found",
			itemId: "2",
			setupMock: func(iu *mock_usecase.MockItemUseCase) {
				iu.EXPECT().
					GetItemDetailByID(2).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: func() *openapi.ItemDetail {
				return nil
			},
		},
		{
			name:   "failure: invalid item ID",
			itemId: "abc",
			setupMock: func(iu *mock_usecase.MockItemUseCase) {
				// No calls expected
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: func() *openapi.ItemDetail {
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemUseCase := mock_usecase.NewMockItemUseCase(ctrl)
			mockCommentUsecase := mock_usecase.NewMockCommentUsecase(ctrl)
			mockFileUseCase := mock_usecase.NewMockFileUseCase(ctrl)
			tc.setupMock(mockItemUseCase)

			h := NewHandler(mockItemUseCase, mockCommentUsecase, mockFileUseCase, nil, nil)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/items/%s", tc.itemId), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			body := strings.TrimSpace(rec.Body.String())

			if tc.expectedCode == http.StatusOK {
				expectedByte, err := tc.expectedBody().MarshalJSON()
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
		setupMock    func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase)
		expectedCode int
		expectedBody func() []openapi.ItemSummary
	}{
		{
			name:  "success: no query",
			query: "",
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				iu.EXPECT().
					SearchItems(domain.ItemSearchQuery{}).
					Return([]*domain.ItemDetail{
						{
							Item: &domain.Item{
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
							Tags:  []*domain.Tag{{Name: "tag-1"}},
							Likes: []*domain.Like{{ItemID: 1, UserID: "user-a"}},
						},
						{
							Item: &domain.Item{
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
							Tags:  []*domain.Tag{},
							Likes: []*domain.Like{{ItemID: 2, UserID: "user-b"}, {ItemID: 2, UserID: "user-c"}},
						},
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() []openapi.ItemSummary {
				equipmentCount := 1
				equipmentIsLiked := false
				equipmentLikeCounts := 1
				equipment := &openapi.ItemSummary{
					Id:          1,
					Name:        "Test Item 1",
					Description: "This is the first test item",
					ImgUrl:      "http://example.com/image1.png",
					Count:       &equipmentCount,
					IsBook:      false,
					IsTrapItem:  true,
					IsLiked:     &equipmentIsLiked,
					LikeCounts:  &equipmentLikeCounts,
					Tags: &[]openapi.Tag{
						{Name: "tag-1"},
					},
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					DeletedAt: nullable.NewNullNullable[time.Time](),
				}

				code := "0987654321"
				bookIsLiked := false
				bookLikeCounts := 2
				book := &openapi.ItemSummary{
					Id:          2,
					Name:        "Test Item 2",
					Description: "This is the second test item",
					ImgUrl:      "http://example.com/image2.png",
					IsBook:      true,
					IsTrapItem:  false,
					IsLiked:     &bookIsLiked,
					LikeCounts:  &bookLikeCounts,
					Code:        &code,
					CreatedAt:   time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nullable.NewNullNullable[time.Time](),
				}
				return []openapi.ItemSummary{*equipment, *book}
			},
		},
		{
			name:  "success: with query",
			query: "?search=test&limit=10",
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				iu.EXPECT().
					SearchItems(domain.ItemSearchQuery{Name: "test", Limit: 10}).
					Return([]*domain.ItemDetail{
						{
							Item: &domain.Item{
								ID:              1,
								Name:            "Test Item",
								Description:     "This is a test item",
								ImgUrl:          "http://example.com/image.png",
								BookDetail:      nil,
								EquipmentDetail: nil,
								CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt:       time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
								DeletedAt:       nil,
							},
							Tags:  []*domain.Tag{},
							Likes: []*domain.Like{},
						},
					}, nil).
					Times(1)
			},

			expectedCode: http.StatusOK,
			expectedBody: func() []openapi.ItemSummary {
				isLiked := false
				likeCounts := 0
				return []openapi.ItemSummary{
					{
						Id:          1,
						Name:        "Test Item",
						Description: "This is a test item",
						ImgUrl:      "http://example.com/image.png",
						IsBook:      false,
						IsLiked:     &isLiked,
						IsTrapItem:  false,
						LikeCounts:  &likeCounts,
						CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
						DeletedAt:   nullable.NewNullNullable[time.Time](),
					},
				}
			},
		},
		{
			name:  "failure: invalid query",
			query: "?limit=-1",
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				iu.EXPECT().
					SearchItems(domain.ItemSearchQuery{Limit: -1}).
					Return(nil, usecase.ErrInvalidSearchQuery).
					Times(1)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: func() []openapi.ItemSummary {
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemUseCase := mock_usecase.NewMockItemUseCase(ctrl)
			mockTagUseCase := mock_usecase.NewMockTagUseCase(ctrl)
			tc.setupMock(mockItemUseCase, mockTagUseCase)

			h := NewHandlerWithTagLike(mockItemUseCase, nil, nil, nil, nil, mockTagUseCase, nil)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/items%s", tc.query), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			body := strings.TrimSpace(rec.Body.String())

			if tc.expectedCode == http.StatusOK {
				expectedByte, err := json.Marshal(tc.expectedBody())
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
		setupMock    func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase)
		expectedCode int
		expectedBody func() []openapi.Item
	}{
		{
			name: "success",
			requestBody: `[
				{
					"name": "New Item",
					"description": "This is a new item",
					"imgUrl": "http://example.com/new_image.png",
					"isBook": true,
					"isTrapItem": false,
					"code": "1234567890"
				}
			]`,
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				iu.EXPECT().
					CreateItems([]*domain.Item{
						{
							Name:        "New Item",
							Description: "This is a new item",
							ImgUrl:      "http://example.com/new_image.png",
							BookDetail: &domain.BookDetail{
								ISBNCode: "1234567890",
							},
							EquipmentDetail: nil,
						},
					}).
					Return([]*domain.Item{
						{
							ID:          1,
							Name:        "New Item",
							Description: "This is a new item",
							ImgUrl:      "http://example.com/new_image.png",
							BookDetail: &domain.BookDetail{
								ISBNCode: "1234567890",
							},
							EquipmentDetail: nil,
							CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
							DeletedAt:       nil,
						},
					}, nil).
					Times(1)
				tu.EXPECT().
					ReplaceByItemID(1, ([]string)(nil)).
					Return(nil).
					Times(1)
				tu.EXPECT().
					GetByItemID(1).
					Return([]*domain.Tag{}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() []openapi.Item {
				code := "1234567890"
				return []openapi.Item{
					{
						Id:          1,
						Name:        "New Item",
						Description: "This is a new item",
						ImgUrl:      "http://example.com/new_image.png",
						IsBook:      true,
						IsTrapItem:  false,
						Code:        &code,
						CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						DeletedAt:   nullable.NewNullNullable[time.Time](),
					},
				}
			},
		},
		{
			name: "success: multiple items",
			requestBody: `[
				{
					"name": "New Item 1",
					"description": "This is the first new item",
					"imgUrl": "http://example.com/new_image1.png",
					"isBook": false,
					"isTrapItem": true,
					"count": 3
				},
				{
					"name": "New Item 2",
					"description": "This is the second new item",
					"imgUrl": "http://example.com/new_image2.png",
					"isBook": true,
					"isTrapItem": false,
					"code": "0987654321"
				}
			]`,
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				gomock.InOrder(
					iu.EXPECT().
						CreateItems([]*domain.Item{
							{
								Name:        "New Item 1",
								Description: "This is the first new item",
								ImgUrl:      "http://example.com/new_image1.png",
								BookDetail:  nil,
								EquipmentDetail: &domain.EquipmentDetail{
									Count:    3,
									CountMax: 3,
								},
							},
							{
								Name:        "New Item 2",
								Description: "This is the second new item",
								ImgUrl:      "http://example.com/new_image2.png",
								BookDetail: &domain.BookDetail{
									ISBNCode: "0987654321",
								},
								EquipmentDetail: nil,
							},
						}).
						Return([]*domain.Item{
							{
								ID:          1,
								Name:        "New Item 1",
								Description: "This is the first new item",
								ImgUrl:      "http://example.com/new_image1.png",
								BookDetail:  nil,
								EquipmentDetail: &domain.EquipmentDetail{
									Count:    3,
									CountMax: 3,
								},
								CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
								DeletedAt: nil,
							},
							{
								ID:          2,
								Name:        "New Item 2",
								Description: "This is the second new item",
								ImgUrl:      "http://example.com/new_image2.png",
								BookDetail: &domain.BookDetail{
									ISBNCode: "0987654321",
								},
								EquipmentDetail: nil,
								CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
								DeletedAt:       nil,
							},
						}, nil).
						Times(1),
					tu.EXPECT().
						ReplaceByItemID(1, ([]string)(nil)).
						Return(nil).
						Times(1),
					tu.EXPECT().
						ReplaceByItemID(2, ([]string)(nil)).
						Return(nil).
						Times(1),
					tu.EXPECT().
						GetByItemID(1).
						Return([]*domain.Tag{}, nil).
						Times(1),
					tu.EXPECT().
						GetByItemID(2).
						Return([]*domain.Tag{}, nil).
						Times(1),
				)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() []openapi.Item {
				equipment := &openapi.Item{
					Id:          1,
					Name:        "New Item 1",
					Description: "This is the first new item",
					ImgUrl:      "http://example.com/new_image1.png",
					IsBook:      false,
					IsTrapItem:  true,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nullable.NewNullNullable[time.Time](),
				}
				_ = equipment.FromItem0(openapi.Item0{
					Count:    3,
					CountMax: 3,
				})

				code := "0987654321"
				book := &openapi.Item{
					Id:          2,
					Name:        "New Item 2",
					Description: "This is the second new item",
					ImgUrl:      "http://example.com/new_image2.png",
					IsBook:      true,
					IsTrapItem:  false,
					Code:        &code,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nullable.NewNullNullable[time.Time](),
				}

				return []openapi.Item{*equipment, *book}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemUseCase := mock_usecase.NewMockItemUseCase(ctrl)
			mockTagUseCase := mock_usecase.NewMockTagUseCase(ctrl)
			tc.setupMock(mockItemUseCase, mockTagUseCase)

			h := NewHandlerWithTagLike(mockItemUseCase, nil, nil, nil, nil, mockTagUseCase, nil)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			body := strings.TrimSpace(rec.Body.String())

			if tc.expectedCode == http.StatusOK {
				expectedByte, err := json.Marshal(tc.expectedBody())
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
		setupMock    func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase)
		expectedCode int
		expectedBody func() *openapi.Item
	}{
		{
			name:   "success",
			itemId: "1",
			requestBody: `{
				"name": "Updated Item",
				"description": "This is an updated item",
				"imgUrl": "http://example.com/updated_image.png",
				"isBook": false,
				"isTrapItem": true,
				"count": 5
			}`,
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				iu.EXPECT().
					UpdateItem(&domain.Item{
						ID:          1,
						Name:        "Updated Item",
						Description: "This is an updated item",
						ImgUrl:      "http://example.com/updated_image.png",
						BookDetail:  nil,
						EquipmentDetail: &domain.EquipmentDetail{
							Count:    5,
							CountMax: 5,
						},
					}).
					Return(&domain.Item{
						ID:          1,
						Name:        "Updated Item",
						Description: "This is an updated item",
						ImgUrl:      "http://example.com/updated_image.png",
						BookDetail:  nil,
						EquipmentDetail: &domain.EquipmentDetail{
							Count:    5,
							CountMax: 5,
						},
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
						DeletedAt: nil,
					}, nil).
					Times(1)
				tu.EXPECT().
					ReplaceByItemID(1, ([]string)(nil)).
					Return(nil).
					Times(1)
				tu.EXPECT().
					GetByItemID(1).
					Return([]*domain.Tag{}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() *openapi.Item {
				item := &openapi.Item{
					Id:          1,
					Name:        "Updated Item",
					Description: "This is an updated item",
					ImgUrl:      "http://example.com/updated_image.png",
					IsBook:      false,
					IsTrapItem:  true,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nullable.NewNullNullable[time.Time](),
				}
				_ = item.FromItem0(openapi.Item0{
					Count:    5,
					CountMax: 5,
				})
				return item
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
				"isTrapItem": true,
				"count": 5
			}`,
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				iu.EXPECT().
					UpdateItem(&domain.Item{
						ID:          2,
						Name:        "Updated Item",
						Description: "This is an updated item",
						ImgUrl:      "http://example.com/updated_image.png",
						BookDetail:  nil,
						EquipmentDetail: &domain.EquipmentDetail{
							Count:    5,
							CountMax: 5,
						},
					}).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: func() *openapi.Item {
				return nil
			},
		},
		{
			name:   "failure: invalid item ID",
			itemId: "abc",
			requestBody: `{
				"name": "Updated Item",
				"description": "This is an updated item",
				"imgUrl": "http://example.com/updated_image.png",
				"isBook": false,
				"isTrapItem": true,
				"count": 5
			}`,
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				// No calls expected
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: func() *openapi.Item {
				return nil
			},
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
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				iu.EXPECT().
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
			expectedBody: func() *openapi.Item {
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemUseCase := mock_usecase.NewMockItemUseCase(ctrl)
			mockTagUseCase := mock_usecase.NewMockTagUseCase(ctrl)
			tc.setupMock(mockItemUseCase, mockTagUseCase)

			h := NewHandlerWithTagLike(mockItemUseCase, nil, nil, nil, nil, mockTagUseCase, nil)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/items/%s", tc.itemId), strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			body := strings.TrimSpace(rec.Body.String())

			if tc.expectedCode == http.StatusOK {
				expectedByte, err := tc.expectedBody().MarshalJSON()
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
		setupMock    func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase)
		expectedCode int
	}{
		{
			name:   "success",
			itemId: "1",
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				iu.EXPECT().
					DeleteItem(1).
					Return(nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "failure: item not found",
			itemId: "2",
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
				iu.EXPECT().
					DeleteItem(2).
					Return(domain.ErrNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "failure: invalid item ID",
			itemId: "abc",
			setupMock: func(iu *mock_usecase.MockItemUseCase, tu *mock_usecase.MockTagUseCase) {
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
			mockTagUseCase := mock_usecase.NewMockTagUseCase(ctrl)
			tc.setupMock(mockItemUseCase, mockTagUseCase)

			h := NewHandlerWithTagLike(mockItemUseCase, nil, nil, nil, nil, mockTagUseCase, nil)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/items/%s", tc.itemId), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
