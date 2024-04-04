package router

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/model"
)

// PostOwners POST /items/:id/owners
func PostOwners(c echo.Context) error {
	ID := c.Param("id")
	me := c.Get("user").(string)

	body := model.PostOwnershipBody{}
	if err := c.Bind(&body); err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusBadRequest, "リクエストが不正です")
	}

	itemID, err := strconv.Atoi(ID)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusBadRequest, "物品のIDが整数ではありません")
	}

	item, err := model.GetItem(itemID)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusNotFound, "物品がみつかりません")
	}

	if item.Equipment != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusBadRequest, "備品に対して所有者を設定することはできません。物品の情報を変更することで所有物の個数を変更してください")
	}

	ownership := model.Ownership{
		ItemID:     itemID,
		UserID:     me,
		Rentalable: body.Rentalable,
		Memo:       body.Memo,
	}

	ownerRes, err := model.RegisterOwnership(ownership)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusInternalServerError, "所有者の情報を保存できませんでした")
	}

	return c.JSON(http.StatusOK, ownerRes)
}

// PatchOwners PUT /items/:id/owners/:ownershipid
func PatchOwners(c echo.Context) error {
	body := model.PostOwnershipBody{}
	if err := c.Bind(&body); err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusBadRequest, "リクエストが不正です")
	}

	ownershipIDRaw := c.Param("ownershipid")
	ownershipID, err := strconv.Atoi(ownershipIDRaw)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusBadRequest, "変更しようとしている所有者情報のIDが不正です")
	}
	ownership, err := model.GetOwnershipByID(ownershipID)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusNotFound, "変更除しようとしている所有者情報が見つかりません")
	}

	me := c.Get("user").(string)
	if ownership.UserID != me {
		c.Logger().Debug("Error: PatchOwners() " + ownership.UserID + " != " + me)
		return c.JSON(http.StatusForbidden, "権限がないため所有者情報を変更できません")
	}

	ID := c.Param("id")
	itemID, err := strconv.Atoi(ID)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusBadRequest, "物品のIDが整数ではありません")
	}

	ownershipNew := model.Ownership{
		GormModel:  model.GormModel{ID: ownershipID},
		ItemID:     itemID,
		UserID:     me,
		Rentalable: body.Rentalable,
		Memo:       body.Memo,
	}

	err = model.PatchOwnership(ownershipNew)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusInternalServerError, "所有者情報を変更できません")
	}

	return c.JSON(http.StatusAccepted, ownershipNew)
}

// DeleteOwners PUT /items/:id/owners/:ownershipid
func DeleteOwners(c echo.Context) error {
	ownershipIDRaw := c.Param("ownershipid")
	me := c.Get("user").(string)

	ownershipID, err := strconv.Atoi(ownershipIDRaw)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusBadRequest, "削除しようとしている所有者情報のIDが不正です")
	}

	ownership, err := model.GetOwnershipByID(ownershipID)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusNotFound, "削除しようとしている所有者情報が見つかりません")
	}

	if ownership.UserID != me {
		c.Logger().Debug("Error: DeleteOwners() " + ownership.UserID + " != " + me)
		return c.JSON(http.StatusForbidden, "権限がないため所有者情報を削除できません")
	}

	err = model.DeleteOwnership(ownershipID)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSON(http.StatusInternalServerError, "所有者情報を削除できません")
	}

	return c.NoContent(http.StatusOK)
}
