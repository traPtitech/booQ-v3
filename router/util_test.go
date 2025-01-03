package router

import (
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
)

var TEST_USER = "s9"

func performMutation(e *echo.Echo, method, path, payload string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}
