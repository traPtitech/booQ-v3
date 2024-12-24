package router

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/model"
)

// GetItems GET /items
func GetItems(c echo.Context) error {
	getItemsBody, err := parseGetItemsParams(c)
	if err != nil {
		return invalidRequest(c, err)
	}

	res, err := model.GetItems(getItemsBody)
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func parseGetItemsParams(c echo.Context) (model.GetItemsBody, error) {
	query := c.QueryParams()

	getItemsBody := model.GetItemsBody{
		UserID:      query.Get("userId"),
		Search:      query.Get("search"),
		Rental:      query.Get("rental"),
		Tags:        query["tag"],
		TagsExclude: query["tag-exclude"],
		SortBy:      query.Get("sortby"),
	}

	if query.Get("limit") != "" {
		limit, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			return model.GetItemsBody{}, err
		}
		getItemsBody.Limit = limit
	}

	if query.Get("offset") != "" {
		offset, err := strconv.Atoi(query.Get("offset"))
		if err != nil {
			return model.GetItemsBody{}, err
		}
		getItemsBody.Offset = offset
	}

	c.Logger().Info(getItemsBody)

	return getItemsBody, nil
}

// PostItems POST /items
func PostItems(c echo.Context) error {
	me, err := getAuthorizedUser(c)
	if err != nil {
		return unauthorizedRequest(c, err)
	}

	items := []model.RequestPostItemsBody{}
	if err := c.Bind(&items); err != nil {
		return invalidRequest(c, err)
	}

	res, err := model.CreateItems(items, me)
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

// GetItem GET /items/:id
func GetItem(c echo.Context) error {
	itemIDRaw := c.Param("id")

	itemID, err := strconv.Atoi(itemIDRaw)
	if err != nil {
		return invalidRequest(c, err)
	}

	res, err := model.GetItem(itemID)
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

// PatchItem PUT /items/:id
func PatchItem(c echo.Context) error {
	itemBody := model.RequestPostItemsBody{}
	err := c.Bind(&itemBody)
	if err != nil {
		return invalidRequest(c, err)
	}

	itemIDRaw := c.Param("id")
	itemID, err := strconv.Atoi(itemIDRaw)
	if err != nil {
		return invalidRequest(c, err)
	}

	res, err := model.PatchItem(itemID, itemBody)
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteItem DELETE /items/:id
func DeleteItem(c echo.Context) error {
	itemIDRaw := c.Param("id")
	itemID, err := strconv.Atoi(itemIDRaw)
	if err != nil {
		return invalidRequest(c, err)
	}

	err = model.DeleteItem(itemID)
	if err != nil {
		return internalServerError(c, err)
	}

	return c.NoContent(http.StatusOK)
}
