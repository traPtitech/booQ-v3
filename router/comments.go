package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostComments POST /items/:id/comments
func PostComments(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
