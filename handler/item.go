package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/usecase"
)

func (h *handler) GetItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	item, err := h.iu.GetItemByID(itemId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to get item: %v", err))
	}

	return ctx.JSON(http.StatusOK, toOpenAPIItem(item))
}

func (h *handler) GetItems(ctx echo.Context, params openapi.GetItemsParams) error {
	query := domain.ItemSearchQuery{}

	if params.UserId != nil {
		query.UserID = *params.UserId
	}
	if params.Search != nil {
		query.Name = *params.Search
	}
	if params.Rental != nil {
		query.BorrowerID = *params.Rental
	}
	if params.Limit != nil {
		query.Limit = *params.Limit
	}
	if params.Offset != nil {
		query.Offset = *params.Offset
	}
	if params.Tag != nil {
		query.Tag = *params.Tag
	}
	if params.TagExclude != nil {
		query.TagExclude = *params.TagExclude
	}
	if params.Sortby != nil {
		query.SortBy = *params.Sortby
	}

	items, err := h.iu.SearchItems(query)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidSearchQuery) {
			return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid search query: %v", err))
		}
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to search items: %v", err))
	}

	openAPIItems := make([]openapi.Item, 0, len(items))
	for _, item := range items {
		openAPIItems = append(openAPIItems, *toOpenAPIItem(item))
	}

	return ctx.JSON(http.StatusOK, openAPIItems)
}

func (h *handler) PostItem(ctx echo.Context) error {
	itemRequest := openapi.PostItemJSONRequestBody{}
	if err := ctx.Bind(&itemRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
	}

	items := make([]*domain.Item, 0, len(itemRequest))
	for _, req := range itemRequest {
		item, err := postRequestToDomainItem(&req)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid item in request body: %v", err))
		}
		items = append(items, item)
	}

	createdItems, err := h.iu.CreateItems(items)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to create items: %v", err))
	}

	openAPIItems := make([]openapi.Item, 0, len(createdItems))
	for _, item := range createdItems {
		openAPIItems = append(openAPIItems, *toOpenAPIItem(item))
	}

	return ctx.JSON(http.StatusOK, openAPIItems)
}

func (h *handler) DeleteItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	err := h.iu.DeleteItem(itemId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to delete item: %v", err))
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *handler) EditItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	itemRequest := openapi.EditItemJSONRequestBody{}
	if err := ctx.Bind(&itemRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
	}

	item, err := postRequestToDomainItem(&itemRequest)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid item in request body: %v", err))
	}
	item.ID = itemId

	updatedItem, err := h.iu.UpdateItem(item)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		if errors.Is(err, usecase.ErrUpdateNotAllowed) {
			return ctx.String(http.StatusBadRequest, fmt.Sprintf("update not allowed: %v", err))
		}
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to update item: %v", err))
	}

	return ctx.JSON(http.StatusOK, toOpenAPIItem(updatedItem))
}
