package router

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/model"
)

// GetItems GET /items
func GetItems(c echo.Context) error {
	getItemsBody, err := getItemsParams(c)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusBadRequest, "リクエストデータの処理に失敗しました")
	}

	res, err := model.GetItems(getItemsBody)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusInternalServerError, "DBの操作に失敗しました")
	}

	return c.JSON(http.StatusOK, res)
}

func getItemsParams(c echo.Context) (model.GetItemsBody, error) {
	params := c.QueryParams()

	getItemsBody := model.GetItemsBody{
		UserID:      params.Get("userId"),
		Search:      params.Get("search"),
		Rental:      params.Get("rental"),
		Tags:        params["tag"],
		TagsExclude: params["tag-exclude"],
		SortBy:      params.Get("sortby"),
	}

	if params.Get("limit") != "" {
		limit, err := strconv.Atoi(params.Get("limit"))
		if err != nil {
			return model.GetItemsBody{}, err
		}
		getItemsBody.Limit = limit
	}

	if params.Get("offset") != "" {
		offset, err := strconv.Atoi(params.Get("offset"))
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
	me := c.Get("user").(string)
	items := []model.RequestPostItemsBody{}
	err := c.Bind(&items)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusBadRequest, "リクエストデータの処理に失敗しました")
	}

	res, err := model.CreateItems(items, me)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusBadRequest, "DBの操作に失敗しました")
	}

	return c.JSON(http.StatusOK, res)
}

// GetItem GET /items/:id
func GetItem(c echo.Context) error {
	itemIDRaw := c.Param("id")

	itemID, err := strconv.Atoi(itemIDRaw)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusBadRequest, "物品のIDが不正です")
	}

	res, err := model.GetItem(itemID)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusInternalServerError, "DBの操作に失敗しました")
	}

	return c.JSON(http.StatusOK, res)
}

// PatchItem PUT /items/:id
func PatchItem(c echo.Context) error {
	itemBody := model.RequestPostItemsBody{}
	err := c.Bind(&itemBody)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusBadRequest, "リクエストデータの処理に失敗しました")
	}

	itemIDRaw := c.Param("id")
	itemID, err := strconv.Atoi(itemIDRaw)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusBadRequest, "物品のIDが不正です")
	}

	res, err := model.PatchItem(itemID, itemBody)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusInternalServerError, "DBの操作に失敗しました")
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteItem DELETE /items/:id
func DeleteItem(c echo.Context) error {
	itemIDRaw := c.Param("id")
	itemID, err := strconv.Atoi(itemIDRaw)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusBadRequest, "物品のIDが不正です")
	}

	err = model.DeleteItem(itemID)
	if err != nil {
		c.Logger().Info(err.Error())
		return c.JSON(http.StatusInternalServerError, "DBの操作に失敗しました")
	}

	return c.NoContent(http.StatusOK)
}
