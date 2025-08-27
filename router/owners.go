package router

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/model"
)

// PostOwners POST /items/:id/owners
func PostOwners(c echo.Context) error {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return invalidRequest(c, err)
	}

	owner := model.OwnershipPayload{}
	if err := c.Bind(&owner); err != nil {
		return invalidRequest(c, err)
	}
	owner.ItemID = itemId

	res, err := model.CreateOwnership(owner)
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

// PatchOwners PUT /items/:id/owners/:ownershipid
func PatchOwners(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// DeleteOwners PUT /items/:id/owners/:ownershipid
func DeleteOwners(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
