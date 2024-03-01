package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetItems GET /items
func GetItems(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// PostItems POST /items
func PostItems(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// GetItem GET /items/:id
func GetItem(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// PutItem PUT /items/:id
func PutItem(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// DeleteItem DELETE /items/:id
func DeleteItem(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
