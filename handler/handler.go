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
	bu usecase.BorrowingUseCase
}

func NewHandler(iu usecase.ItemUseCase, fu usecase.FileUseCase, ou usecase.OwnershipUseCase, bu usecase.BorrowingUseCase) openapi.ServerInterface {
	return &handler{
		fu: fu,
		iu: iu,
		ou: ou,
		bu: bu,
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
