package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
)

func (h *handler) PostComment(ctx echo.Context, itemId openapi.ItemIdInPath) error {

	var req openapi.PostComment
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// TODO: ミドルウェアからユーザーidを取得するように
	userId := "test-user-id"

	comment, err := h.cu.CreateComment(itemId, userId, req.Text)
	if err != nil {
		if err == domain.ErrItemNotFound {
			// Itemがないとき
			return ctx.JSON(http.StatusNotFound, "item not found")
		} else if err == domain.ErrCommentTextEmpty {
			// 投稿されたコメントが空の時
			return ctx.JSON(http.StatusBadRequest, "comment text is empty")
		}

		return ctx.JSON(http.StatusInternalServerError, err)
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
