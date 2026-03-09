package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/oapi-codegen/nullable"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
)

func toOpenAPIItem(domainItem *domain.Item) (*openapi.Item, error) {
	if domainItem == nil {
		return nil, errors.New("domain item is nil")
	}

	var deletedAt nullable.Nullable[time.Time]
	if domainItem.DeletedAt != nil {
		deletedAt = nullable.NewNullableWithValue(*domainItem.DeletedAt)
	} else {
		deletedAt = nullable.NewNullNullable[time.Time]()
	}

	// TODO: add tags, comments, likes, transactions, etc.
	item := &openapi.Item{
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

	if domainItem.BookDetail != nil {
		item.Code = &domainItem.BookDetail.ISBNCode
	}
	if domainItem.EquipmentDetail != nil {
		i := openapi.Item0{
			Count:    domainItem.EquipmentDetail.Count,
			CountMax: domainItem.EquipmentDetail.CountMax,
		}
		err := item.FromItem0(i)
		if err != nil {
			return nil, fmt.Errorf("failed to convert equipment detail: %w", err)
		}
	}

	return item, nil
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
			return nil, errors.New("count is required for trap items")
		}
		item.EquipmentDetail = &domain.EquipmentDetail{
			Count:    *request.Count,
			CountMax: *request.Count, // 作ったときはすべての備品が貸出可能
		}
	}

	return item, nil
}

func toOpenAPIOwnership(d *domain.Ownership) openapi.Ownership {
	return openapi.Ownership{
		Id:         &d.ID,
		ItemId:     &d.ItemID,
		UserId:     d.UserID,
		Rentalable: d.Rentable,
		Memo:       d.Memo,
	}
}
