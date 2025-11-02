package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/handler/openapi"
)

func (h *Handler) GetItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	item, err := h.iu.GetItemByID(itemId)
	if err != nil {
		// TODO: NotFoundなどのエラーも考慮する
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, item)
}
