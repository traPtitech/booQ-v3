package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/usecase"
)

type handler struct {
	iu usecase.ItemUseCase
}

func NewHandler(iu usecase.ItemUseCase) openapi.ServerInterface {
	return &handler{
		iu: iu,
	}
}

func (h *handler) PostFile(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) GetFile(ctx echo.Context, fileId openapi.FileIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) GetItems(ctx echo.Context, params openapi.GetItemsParams) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) PostItem(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) DeleteItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) EditItem(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) PostBorrowEquipment(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) PostBorrowEquipmentReturn(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) PostComment(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) RemoveLike(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) AddLike(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) PostItemOwners(ctx echo.Context, itemId openapi.ItemIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) DeleteItemOwners(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) EditItemOwners(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) PostBorrow(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) GetBorrowingById(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) PostBorrowReply(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) PostReturn(ctx echo.Context, itemId openapi.ItemIdInPath, ownershipId openapi.OwnershipIdInPath, borrowingId openapi.BorrowingIdInPath) error {
	//TODO implement me
	panic("implement me")
}
