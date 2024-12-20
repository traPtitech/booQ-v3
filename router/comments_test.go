package router

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/model"
)

func TestPostComment(t *testing.T) {
	model.PrepareTestDatabase()

	e := echo.New()
	SetupRouting(e, CreateUserProvider(TEST_USER))

	cases := []struct {
		name     string
		payload  string
		expected int
	}{
		{
			name:     "正常系",
			payload:  `{"text":"テストコメント"}`,
			expected: http.StatusCreated,
		},
		{
			name:     "異常系: 空文字列",
			payload:  `{"text":""}`,
			expected: http.StatusBadRequest,
		},
		{
			name:     "異常系: パラメータ不足",
			payload:  `{}`,
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rec := performMutation(e, "POST", "/api/items/1/comments", tc.payload)
			assert.Equal(tc.expected, rec.Code)
		})
	}
}
