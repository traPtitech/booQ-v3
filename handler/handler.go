package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/usecase"
)

type Handler struct {
	iu usecase.ItemUseCase
}

func (h *Handler) PostFile(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GetFile(ctx echo.Context, fileId openapi.FileIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GetItems(ctx echo.Context, params openapi.GetItemsParams) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PostItem(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeleteItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GetItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) EditItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PostBorrowEquipment(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PostBorrowEquipmentReturn(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PostComment(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) RemoveLike(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) AddLike(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PostItemOwners(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeleteItemOwners(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) EditItemOwners(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PostBorrow(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GetBorrowingById(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PostBorrowReply(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PostReturn(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	//TODO implement me
	panic("implement me")
}
