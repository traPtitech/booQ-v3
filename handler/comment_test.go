package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	mock_usecase "github.com/traPtitech/booQ-v3/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestHandler_PostComment(t *testing.T) {
	type fields struct {
		itemUsecase    *mock_usecase.MockItemUseCase
		commentUsecase *mock_usecase.MockCommentUsecase
		fileUsecase    *mock_usecase.MockFileUseCase
	}

	testCases := []struct {
		name         string
		itemId       int
		requestBody  string
		setupMock    func(f *fields)
		expectedCode int
		expectedBody string
	}{
		{
			name:        "success",
			itemId:      1,
			requestBody: `{"text":"test comment"}`,
			setupMock: func(f *fields) {
				f.commentUsecase.EXPECT().CreateComment(1, "test-user-id", "test comment").Return(&domain.Comment{
					ID:        10,
					ItemID:    1,
					UserID:    "test-user-id",
					Text:      "test comment",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: "test comment",
		},
		{
			name:         "bad request: invalid body",
			itemId:       1,
			requestBody:  `invalid json`,
			setupMock:    func(f *fields) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "not found",
			itemId:      999,
			requestBody: `{"text":"test comment"}`,
			setupMock: func(f *fields) {
				f.commentUsecase.EXPECT().CreateComment(999, "test-user-id", "test comment").Return(nil, domain.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "not found",
		},
		{
			name:        "bad request: empty comment",
			itemId:      1,
			requestBody: `{"text":""}`,
			setupMock: func(f *fields) {
				f.commentUsecase.EXPECT().CreateComment(1, "test-user-id", "").Return(nil, domain.ErrCommentTextEmpty)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "comment text cannot be empty",
		},
		{
			name:        "internal error",
			itemId:      1,
			requestBody: `{"text":"test comment"}`,
			setupMock: func(f *fields) {
				f.commentUsecase.EXPECT().CreateComment(1, "test-user-id", "test comment").Return(nil, errors.New("unexpected error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "unexpected error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				itemUsecase:    mock_usecase.NewMockItemUseCase(ctrl),
				commentUsecase: mock_usecase.NewMockCommentUsecase(ctrl),
				fileUsecase:    mock_usecase.NewMockFileUseCase(ctrl),
			}
			tc.setupMock(f)

			h := NewHandler(f.itemUsecase, f.commentUsecase, f.fileUsecase, nil, nil)
			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req := httptest.NewRequest(http.MethodPost, "/items/"+strconv.Itoa(tc.itemId)+"/comments", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/items/:itemId/comments")

			c.SetParamNames("itemId")
			c.SetParamValues(strconv.Itoa(tc.itemId))

			_ = h.PostComment(c, openapi.ItemIdInPath(tc.itemId))

			if tc.expectedCode == http.StatusCreated {
				assert.Equal(t, tc.expectedCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
			} else {
				assert.Equal(t, tc.expectedCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
			}
		})
	}
}
