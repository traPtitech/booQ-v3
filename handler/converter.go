package handler

import (
	"time"

	"github.com/oapi-codegen/nullable"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
)

func toOpenAPIItem(domainItem *domain.Item) *openapi.Item {
	if domainItem == nil {
		return nil
	}

	var deletedAt nullable.Nullable[time.Time]
	if domainItem.DeletedAt != nil {
		deletedAt = nullable.NewNullableWithValue(*domainItem.DeletedAt)
	} else {
		deletedAt = nullable.NewNullNullable[time.Time]()
	}

	// TODO: add tags, comments, likes, transactions, etc.
	return &openapi.Item{
		Id:          domainItem.ID,
		Name:        domainItem.Name,
		Description: domainItem.Description,
		ImgUrl:      domainItem.ImgUrl,
		IsBook:      domainItem.BookDetail != nil,
		IsTrapItem:  domainItem.EquipmentDetail != nil,
		CreatedAt:   domainItem.CreatedAt,
		UpdatedAt:   domainItem.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}
