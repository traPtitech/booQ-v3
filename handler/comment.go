package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/middleware"
)

func (h *handler) PostComment(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	var req openapi.PostComment
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request body")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	comment, err := h.cu.CreateComment(itemId, userID, req.Text)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, "not found")
		} else if errors.Is(err, domain.ErrCommentTextEmpty) {
			return ctx.JSON(http.StatusBadRequest, "comment text cannot be empty")
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	res := openapi.Comment{
		Id:        &comment.ID,
		ItemId:    &comment.ItemID,
		UserId:    &comment.UserID,
		Text:      comment.Text,
		CreatedAt: &comment.CreatedAt,
		UpdatedAt: &comment.UpdatedAt,
	}

	return ctx.JSON(http.StatusCreated, res)
}
