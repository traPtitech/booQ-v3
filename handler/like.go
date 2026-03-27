package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/middleware"
	"github.com/traPtitech/booQ-v3/usecase"
)

func (h *handler) RemoveLike(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	if h.lu == nil {
		return ctx.JSON(http.StatusInternalServerError, "like usecase is not configured")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	err := h.lu.RemoveLike(itemId, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		if errors.Is(err, usecase.ErrNotLiked) {
			return ctx.JSON(http.StatusBadRequest, "item is not liked")
		}
		return ctx.JSON(http.StatusInternalServerError, "failed to remove like")
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *handler) AddLike(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	if h.lu == nil {
		return ctx.JSON(http.StatusInternalServerError, "like usecase is not configured")
	}

	userID, ok := middleware.GetUserID(ctx.Request().Context())
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "user ID not found in context")
	}

	err := h.lu.AddLike(itemId, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		if errors.Is(err, usecase.ErrAlreadyLiked) {
			return ctx.JSON(http.StatusBadRequest, "item is already liked")
		}
		return ctx.JSON(http.StatusInternalServerError, "failed to add like")
	}

	return ctx.NoContent(http.StatusCreated)
}
