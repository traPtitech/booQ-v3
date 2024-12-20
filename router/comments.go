package router

import (
	"fmt"
	"net/http"
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
	itemIDStr := c.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		return invalidRequest(c, err)
	}

	me, err := getAuthorizedUser(c)
	if err != nil {
		return unauthorizedRequest(c, err)
	}

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

	return c.JSON(http.StatusCreated, PostCommentResponse{ID: comment.ID})
}
