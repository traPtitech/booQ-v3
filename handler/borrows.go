package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
)

// POST /items/:itemId/borrowing/equipment
func (h *handler) PostBorrowEquipment(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	var reqBody domain.BorrowRequestEquipment
	if err := ctx.Bind(&reqBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid request body")
	}

	if reqBody.DueDate == "" {
		return ctx.JSON(http.StatusBadRequest, "dueDate is required")
	}

	userID := 1

	res, err := h.bu.BorrowEquipment(int(itemId), userID, reqBody)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrItemNotFound):
			return ctx.NoContent(http.StatusNotFound)
		case errors.Is(err, domain.ErrNotEnoughStock):
			return ctx.JSON(http.StatusConflict, map[string]string{"message": "not enough stock"})
		case errors.Is(err, domain.ErrItemNotEquipment):
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "item is not equipment"})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
	}

	return ctx.JSON(http.StatusCreated, res)
}

// POST /items/:itemId/borrowing/equipment/return
func (h *handler) PostBorrowEquipmentReturn(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	var reqBody domain.BorrowReturn
	if err := ctx.Bind(&reqBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid request body")
	}

	if reqBody.Text == "" {
		return ctx.JSON(http.StatusBadRequest, "text is required")
	}

	userID := 1

	_, err := h.bu.ReturnEquipment(int(itemId), userID, reqBody)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBorrowingNotFound):
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "no active borrowing found"})
		case errors.Is(err, domain.ErrItemNotFound):
			return ctx.NoContent(http.StatusNotFound)
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
	}

	return ctx.JSON(http.StatusCreated, reqBody)
}
