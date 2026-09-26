package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/traPtitech/booQ-v3/internal/domain"
	"github.com/traPtitech/booQ-v3/internal/handler/openapi"
	"github.com/traPtitech/booQ-v3/internal/middleware"
	"github.com/traPtitech/booQ-v3/internal/usecase"
)

func (h *handler) PostBorrow(ctx echo.Context, _ openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
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

	post, err := h.bu.PostRequest(userID, ownershipId, purpose, date, request.BorrowInClubRoom)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidDueDate) {
			return ctx.JSON(http.StatusBadRequest, "due date must be in the future")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, "item or ownership not found")
		}
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

func (h *handler) GetBorrowingById(ctx echo.Context, _ openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	borrowing, err := h.bu.GetRequest(userID, ownershipId, borrowingId)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			return ctx.JSON(http.StatusForbidden, "you cannot access this borrow request")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, "borrow request not found")
		}
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

func (h *handler) PostBorrowReply(ctx echo.Context, _ openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	request := openapi.PostBorrowReplyJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request body")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	reply, err := h.bu.ReplyRequest(userID, ownershipId, borrowingId, request.Answer, request.Comment)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			return ctx.JSON(http.StatusForbidden, "you cannot reply to this borrow request")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, "borrow request not found")
		}
		return ctx.JSON(http.StatusInternalServerError, "failed to reply to borrow request")
	}

	res := openapi.BorrowReply{
		Answer:  request.Answer,
		Comment: reply.Message,
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *handler) PostReturn(ctx echo.Context, _ openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	request := openapi.PostReturnJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request body")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	err := h.bu.ReturnItem(userID, ownershipId, borrowingId, request.Text)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			return ctx.JSON(http.StatusForbidden, "you cannot return this borrow request")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, "borrow request not found")
		}
		return ctx.JSON(http.StatusInternalServerError, "failed to return item")
	}

	return ctx.NoContent(http.StatusOK)
}

// POST /items/:itemId/borrowing/equipment
func (h *handler) PostBorrowEquipment(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	var request openapi.PostBorrowEquipmentJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request body")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	if request.DueDate.IsZero() {
		return ctx.JSON(http.StatusBadRequest, "dueDate is required")
	}
	y, m, d := request.DueDate.Date()
	dueDate := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
	purpose := ""
	if request.Propose != nil {
		purpose = *request.Propose
	}

	count := 1
	if request.Count != nil {
		count = *request.Count
	}

	transaction, err := h.bu.BorrowEquipment(itemId, userID, purpose, count, dueDate)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return ctx.NoContent(http.StatusNotFound)
		case errors.Is(err, usecase.ErrInvalidDueDate), errors.Is(err, usecase.ErrItemNotEquipment):
			return ctx.JSON(http.StatusBadRequest, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, "failed to borrow equipment")
		}
	}

	return ctx.JSON(http.StatusCreated, openapi.BorrowRequestEquipment{
		Propose: &transaction.Purpose, Count: &transaction.Count,
		DueDate: openapi_types.Date{Time: transaction.DueDate}, BorrowInClubRoom: request.BorrowInClubRoom,
	})
}

// POST /items/:itemId/borrowing/equipment/:borrowingId/return
func (h *handler) PostBorrowEquipmentReturn(ctx echo.Context, itemId openapi.ItemIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	var request openapi.PostBorrowEquipmentReturnJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid request body")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	_, err := h.bu.ReturnEquipment(itemId, borrowingId, userID, request.Text)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return ctx.JSON(http.StatusNotFound, "borrowing not found")
		case errors.Is(err, usecase.ErrForbidden):
			return ctx.JSON(http.StatusForbidden, "you cannot return this borrowing")
		case errors.Is(err, domain.ErrInvalidTransactionStatus):
			return ctx.JSON(http.StatusBadRequest, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, "failed to return equipment")
		}
	}
	
	return ctx.JSON(http.StatusCreated, openapi.BorrowReturn(request))
}
