package handler

import (
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
)

func toOpenAPIItem(domainItem *domain.Item) *openapi.Item {
	if domainItem == nil {
		return nil
	}

	// TODO: add tags, comments, likes, transactions, etc.
	return &openapi.Item{
		Id:          &domainItem.ID,
		Name:        domainItem.Name,
		Description: domainItem.Description,
		ImgUrl:      domainItem.ImgUrl,
		IsBook:      domainItem.BookDetail != nil,
		IsTrapItem:  domainItem.EquipmentDetail != nil,
	}
}
