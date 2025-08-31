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
		return parseModelError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

// PatchOwners PUT /items/:id/owners/:ownershipid
func PatchOwners(c echo.Context) error {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return invalidRequest(c, err)
	}

	ownershipId, err := strconv.Atoi(c.Param("ownershipid"))
	if err != nil {
		return invalidRequest(c, err)
	}

	owner := model.OwnershipPayload{}
	if err := c.Bind(&owner); err != nil {
		return invalidRequest(c, err)
	}
	owner.ItemID = itemId

	res, err := model.UpdateOwnership(ownershipId, owner)
	if err != nil {
		return parseModelError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteOwners PUT /items/:id/owners/:ownershipid
func DeleteOwners(c echo.Context) error {
	ownershipId, err := strconv.Atoi(c.Param("ownershipid"))
	if err != nil {
		return invalidRequest(c, err)
	}

	me, err := getAuthorizedUser(c)
	if err != nil {
		return unauthorizedRequest(c, err)
	}

	err = model.DeleteOwnership(ownershipId, me)
	if err != nil {
		return parseModelError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
