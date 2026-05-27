package handler

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
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

// TestHandler_PostFile は POST /files のテスト
func TestHandler_PostFile(t *testing.T) {
	// テストケースの定義（テーブル駆動テスト）
	testCases := []struct {
		name         string
		setupRequest func() (*http.Request, error)
		setupMock    func(f *mock_usecase.MockFileUseCase)
		expectedCode int
		expectedBody string
	}{
		{
			name: "success", // 正常系のテスト
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)
				h.Set("Content-Type", "image/jpeg")
				part, err := writer.CreatePart(h)
				if err != nil {
					return nil, err
				}
				_, err = part.Write([]byte("fake image data"))
				if err != nil {
					return nil, err
				}
				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/files", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			setupMock: func(f *mock_usecase.MockFileUseCase) {
				// Upload メソッドが呼ばれることを期待
				f.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&domain.File{
						ID:        1,
						Name:      "abc123.jpg",
						MimeType:  "image/jpeg",
						CreatedAt: time.Now(),
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusCreated,
			expectedBody: `{"id":1,"url":"/api/files/1"}`,
		},
		{
			name: "failure: file too large", // ファイルサイズ超過のテスト
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)
				h.Set("Content-Type", "image/jpeg")
				part, err := writer.CreatePart(h)
				if err != nil {
					return nil, err
				}
				_, err = part.Write([]byte("fake image data"))
				if err != nil {
					return nil, err
				}
				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/files", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			setupMock: func(f *mock_usecase.MockFileUseCase) {
				f.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, domain.ErrFileTooLarge).
					Times(1)
			},
			expectedCode: http.StatusBadRequest,               // 400 Bad Request を期待
			expectedBody: `"file too large: max size is 3MB"`, // エラーメッセージを期待
		},
		{
			name: "failure: invalid file type", // 不正なファイル形式のテスト
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", "test.gif") // GIF（許可されていない）
				if err != nil {
					return nil, err
				}
				_, err = part.Write([]byte("fake image data"))
				if err != nil {
					return nil, err
				}
				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/files", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			setupMock: func(f *mock_usecase.MockFileUseCase) {
				// invalid file typeはUploadが呼ばれないためmock不要
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"invalid file type: only JPEG and PNG are allowed"`,
		},
		{
			name: "failure: no file", // ファイルがないテスト
			setupRequest: func() (*http.Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/files", nil)
				return req, nil
			},
			setupMock: func(f *mock_usecase.MockFileUseCase) {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"file is required"`,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// モックを作成
			mockItemUseCase := mock_usecase.NewMockItemUseCase(ctrl)
			mockFileUseCase := mock_usecase.NewMockFileUseCase(ctrl)
			tc.setupMock(mockFileUseCase)

			h := NewHandler(mockItemUseCase, mockFileUseCase, nil, nil)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			req, err := tc.setupRequest()
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
			body := strings.TrimSpace(rec.Body.String())
			assert.Equal(t, tc.expectedBody, body)
		})
	}
}

// TestHandler_GetFile は GET /files/{fileId} のテスト
func TestHandler_GetFile(t *testing.T) {
	testCases := []struct {
		name         string
		fileId       string
		setupMock    func(f *mock_usecase.MockFileUseCase)
		expectedCode int
		expectedType string
		expectedBody string
	}{
		{
			name:   "success",
			fileId: "1",
			setupMock: func(f *mock_usecase.MockFileUseCase) {
				f.EXPECT().
					GetFile(1).
					Return(
						io.NopCloser(bytes.NewReader([]byte("fake image data"))),
						&domain.File{
							ID:       1,
							Name:     "abc123.jpg",
							MimeType: "image/jpeg",
						},
						nil,
					).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedType: "image/jpeg",
			expectedBody: "fake image data",
		},
		{
			name:   "failure: file not found",
			fileId: "999", // 存在しない ID
			setupMock: func(f *mock_usecase.MockFileUseCase) {
				f.EXPECT().
					GetFile(999).
					Return(nil, nil, domain.ErrNotFound).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
			expectedType: "",
			expectedBody: "",
		},
		{
			name:   "failure: invalid file ID",
			fileId: "abc",
			setupMock: func(f *mock_usecase.MockFileUseCase) {
			},
			expectedCode: http.StatusBadRequest,
			expectedType: "",
			expectedBody: "",
		},
		{
			name:   "failure: internal error",
			fileId: "1",
			setupMock: func(f *mock_usecase.MockFileUseCase) {
				f.EXPECT().
					GetFile(1).
					Return(nil, nil, errors.New("storage error")). // 内部エラー
					Times(1)
			},
			expectedCode: http.StatusInternalServerError,
			expectedType: "",
			expectedBody: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemUseCase := mock_usecase.NewMockItemUseCase(ctrl)
			mockFileUseCase := mock_usecase.NewMockFileUseCase(ctrl)
			tc.setupMock(mockFileUseCase)

			h := NewHandler(mockItemUseCase, mockFileUseCase, nil, nil)

			e := echo.New()
			openapi.RegisterHandlers(e, h)

			// GET /files/{fileId} のリクエストを作成
			req := httptest.NewRequest(http.MethodGet, "/files/"+tc.fileId, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusOK {
				assert.Equal(t, tc.expectedType, rec.Header().Get("Content-Type"))
				assert.Equal(t, tc.expectedBody, rec.Body.String())
			}
		})
	}
}
