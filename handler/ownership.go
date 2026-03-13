package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/middleware"
	"github.com/traPtitech/booQ-v3/usecase"
)

func (h *handler) PostItemOwners(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	req := openapi.PostItemOwnersJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}
	if userID != req.UserId {
		return ctx.JSON(http.StatusForbidden, "you cannot create ownership for another user")
	}

	ownership := &domain.Ownership{
		ItemID:   itemId,
		UserID:   req.UserId,
		Rentable: req.Rentalable,
		Memo:     req.Memo,
	}

	created, err := h.ou.CreateOwnership(ownership)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to create ownership: %v", err))
	}

	return ctx.JSON(http.StatusCreated, toOpenAPIOwnership(created))
}

func (h *handler) EditItemOwners(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
	req := openapi.EditItemOwnersJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	ownership := &domain.Ownership{
		ID:       ownershipId,
		ItemID:   itemId,
		UserID:   req.UserId,
		Rentable: req.Rentalable,
		Memo:     req.Memo,
	}

	updated, err := h.ou.UpdateOwnership(ownership, itemId, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		} else if errors.Is(err, usecase.ErrForbidden) {
			return ctx.JSON(http.StatusForbidden, "you cannot update this ownership")
		}
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to update ownership: %v", err))
	}

	return ctx.JSON(http.StatusOK, toOpenAPIOwnership(updated))
}

func (h *handler) DeleteItemOwners(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	err := h.ou.DeleteOwnership(ownershipId, itemId, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		} else if errors.Is(err, usecase.ErrForbidden) {
			return ctx.JSON(http.StatusForbidden, "you cannot delete this ownership")
		}
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to delete ownership: %v", err))
	}

	return ctx.NoContent(http.StatusOK)
}
