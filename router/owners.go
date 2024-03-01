package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// DeleteItem DELETE /items/:id
func DeleteItem(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// PostOwners POST /items/:id/owners
func PostOwners(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// PatchOwners PUT /items/:id/owners/:ownershipid
func PatchOwners(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// DeleteOwners PUT /items/:id/owners/:ownershipid
func DeleteOwners(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
