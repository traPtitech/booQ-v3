package router

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// PostOwners POST /items/:id/owners
func PostOwners(c echo.Context) error {
	ID := c.Param("id")
	me := c.Get("user").(string)

	itemID, err := strconv.Atoi(ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

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
