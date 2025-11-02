package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
)

func (h *Handler) GetItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	item, err := h.iu.GetItemByID(itemId)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to get item: %v", err))
	}

	return ctx.JSON(http.StatusOK, toOpenAPIItem(item))
}
