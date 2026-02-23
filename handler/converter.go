package handler

import (
	"errors"
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

func postRequestToDomainItem(request *openapi.ItemPostRequest) (*domain.Item, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}

	item := &domain.Item{
		Name:        request.Name,
		Description: request.Description,
		ImgUrl:      request.ImgUrl,
	}

	if request.IsBook {
		if request.Code == nil {
			return nil, errors.New("code is required for book items")
		}
		item.BookDetail = &domain.BookDetail{
			ISBNCode: *request.Code,
		}
	}
	if request.IsTrapItem {
		if request.Count == nil {
			return nil, errors.New("count and countMax are required for trap items")
		}
		item.EquipmentDetail = &domain.EquipmentDetail{
			Count:    *request.Count,
			CountMax: *request.Count, // 作ったときはすべての備品が貸出可能
		}
	}

	return item, nil
}
