package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	for _, tc := range []struct {
		name       string
		headerSet  bool
		userID     string
		wantStatus int
		wantCalled bool
	}{
		{
			name:       "user ID from header",
			headerSet:  true,
			userID:     "user-123",
			wantStatus: http.StatusNoContent,
			wantCalled: true,
		},
		{
			name:       "missing header",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "empty header",
			headerSet:  true,
			wantStatus: http.StatusUnauthorized,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			e.Use(AuthMiddleware)
			called := false
			e.GET("/", func(c echo.Context) error {
				called = true
				userID, ok := GetUserID(c.Request().Context())
				assert.True(t, ok)
				assert.Equal(t, tc.userID, userID)
				return c.NoContent(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.headerSet {
				req.Header.Set("X-Forwarded-User", tc.userID)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantCalled, called)
		})
	}
}
