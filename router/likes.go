package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostLikes POST /items/:id/likes
func PostLikes(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// PostLikes POST /items/:id/likes
func DeleteLikes(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
