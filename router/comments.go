package router

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/model"
)

type PostCommentBody struct {
	Text string `json:"text"`
}

type PostCommentResponse struct {
	ID int `json:"id"`
}

// PostComment POST /items/:id/comments
func PostComment(c echo.Context) error {
	itemIDRaw := c.Param("id")
	itemID, err := strconv.Atoi(itemIDRaw)
	if err != nil {
		return invalidRequest(c, err)
	}

	me := getAuthorizedUser(c)

	var body PostCommentBody
	if err := c.Bind(&body); err != nil {
		return invalidRequest(c, err)
	}
	if body.Text == "" {
		return invalidRequest(c, fmt.Errorf("text is empty"))
	}

	payload := model.CreateCommentPayload{
		ItemID:  itemID,
		UserID:  me,
		Comment: body.Text,
	}
	comment, err := model.CreateComment(&payload)
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(201, PostCommentResponse{ID: comment.ID})
}
