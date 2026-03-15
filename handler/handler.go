package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/usecase"
)

type handler struct {
	fu usecase.FileUseCase
	iu usecase.ItemUseCase
	ou usecase.OwnershipUseCase
}

func NewHandler(iu usecase.ItemUseCase, fu usecase.FileUseCase, ou usecase.OwnershipUseCase) openapi.ServerInterface {
	return &handler{
		fu: fu,
		iu: iu,
		ou: ou,
	}
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
