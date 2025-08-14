package router

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/booQ-v3/model"
)

var testJpeg = `/9j/4AAQSkZJRgABAQIAOAA4AAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////2wBDAf//////////////////////////////////////////////////////////////////////////////////////wAARCAABAAEDAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwBKBH//2Q`

func TestPostFile(t *testing.T) {
	t.Parallel()
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider(TEST_USER))

	t.Run("fail: no form", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		req := httptest.NewRequest(echo.POST, "/api/files", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusBadRequest, rec.Code)
	})

	cases := []struct {
		name          string
		formWriteFunc func(writer *multipart.Writer)
		expected      int
	}{
		{
			name: "fail: invalid file type",
			formWriteFunc: func(writer *multipart.Writer) {
				defer writer.Close()
				part, err := writer.CreateFormFile("file", "test.txt")
				if err != nil {
					t.Error(err)
					return
				}
				_, err = part.Write([]byte("test"))
				if err != nil {
					t.Error(err)
					return
				}
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "fail: bad image",
			formWriteFunc: func(writer *multipart.Writer) {
				defer writer.Close()
				h := textproto.MIMEHeader{}
				h.Set(echo.HeaderContentDisposition, `form-data; name="file"; filename="test.jpg"`)
				h.Set(echo.HeaderContentType, "image/jpeg")
				part, err := writer.CreatePart(h)
				if err != nil {
					t.Error(err)
					return
				}
				_, err = part.Write([]byte("test text file"))
				if err != nil {
					t.Error(err)
					return
				}
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "success",
			formWriteFunc: func(writer *multipart.Writer) {
				defer writer.Close()
				h := textproto.MIMEHeader{}
				h.Set(echo.HeaderContentDisposition, `form-data; name="file"; filename="test.jpg"`)
				h.Set(echo.HeaderContentType, "image/jpeg")
				part, err := writer.CreatePart(h)
				if err != nil {
					t.Error(err)
					return
				}
				img, _ := base64.RawStdEncoding.DecodeString(testJpeg)
				_, err = part.Write(img)
				if err != nil {
					t.Error(err)
					return
				}
			},
			expected: http.StatusCreated,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			pr, pw := io.Pipe()
			writer := multipart.NewWriter(pw)
			go c.formWriteFunc(writer)

			req := httptest.NewRequest(echo.POST, "/api/files", pr)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(c.expected, rec.Code)
		})
	}
}

func TestGetFile(t *testing.T) {
	t.Parallel()
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider(TEST_USER))

	t.Run("fail: invalid id paramater", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		req := httptest.NewRequest(echo.GET, "/api/files/a", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusNotFound, rec.Code)
	})

	t.Run("fail: file not found", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		req := httptest.NewRequest(echo.GET, "/api/files/99999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusNotFound, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		img, _ := base64.RawStdEncoding.DecodeString(testJpeg)
		f, err := model.CreateFile("testuser", bytes.NewReader(img), "jpg")
		require.NoError(t, err)

		req := httptest.NewRequest(echo.GET, fmt.Sprintf("/api/files/%d", f.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(http.StatusOK, rec.Code)
		assert.EqualValues(img, rec.Body.Bytes())
	})
}
