package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostFile POST /files
func PostFile(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// GetFile GET /files/:id
func GetFile(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
