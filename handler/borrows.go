package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/middleware"
)

func (h *handler) PostBorrow(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
	request := openapi.PostBorrowJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request body")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	y, m, d := request.DueDate.Date()
	date := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)

	purpose := ""
	if request.Propose != nil {
		purpose = *request.Propose
	}

	post, err := h.bu.PostRequest(itemId, userID, ownershipId, purpose, date, request.BorrowInClubRoom)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "failed to create borrow request")
	}

	res := openapi.BorrowRequest{
		Propose: &post.Purpose,
		DueDate: openapi_types.Date{
			Time: post.DueDate,
		},
		BorrowInClubRoom: post.BorrowInClubRoom,
	}

	return ctx.JSON(http.StatusCreated, res)
}

func (h *handler) GetBorrowingById(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	borrowing, err := h.bu.GetRequest(itemId, userID, ownershipId, borrowingId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "failed to get borrow request")
	}

	res := openapi.Borrowing{
		Id:               borrowing.ID,
		Propose:          &borrowing.Purpose,
		DueDate:          openapi_types.Date{Time: borrowing.DueDate},
		BorrowInClubRoom: borrowing.BorrowInClubRoom,
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *handler) PostBorrowReply(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	request := openapi.PostBorrowReplyJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request body")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	reply, err := h.bu.ReplyRequest(itemId, userID, ownershipId, borrowingId, request.Answer, request.Comment)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "failed to reply to borrow request")
	}

	res := openapi.BorrowReply{
		Answer:  request.Answer,
		Comment: reply.Message,
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *handler) PostReturn(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	request := openapi.PostReturnJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request body")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	err := h.bu.ReturnItem(itemId, userID, ownershipId, borrowingId, request.Text)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "failed to return item")
	}

	return ctx.NoContent(http.StatusOK)
}
