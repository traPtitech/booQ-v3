package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/oapi-codegen/nullable"
	openapi_types "github.com/oapi-codegen/runtime/types"
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

func toOpenAPIItemSummary(detail *domain.ItemDetail) (*openapi.ItemSummary, error) {
	if detail == nil || detail.Item == nil {
		return nil, errors.New("domain item detail is nil")
	}

	item, err := toOpenAPIItem(detail.Item)
	if err != nil {
		return nil, err
	}

	summary := &openapi.ItemSummary{
		Code:        item.Code,
		Comments:    item.Comments,
		CreatedAt:   item.CreatedAt,
		DeletedAt:   item.DeletedAt,
		Description: item.Description,
		Id:          item.Id,
		ImgUrl:      item.ImgUrl,
		IsBook:      item.IsBook,
		IsTrapItem:  item.IsTrapItem,
		Name:        item.Name,
		Tags:        toOpenAPITags(detail.Tags),
		UpdatedAt:   item.UpdatedAt,
	}
	likeCounts := len(detail.Likes)
	summary.LikeCounts = &likeCounts
	if detail.Item.EquipmentDetail != nil {
		count := detail.Item.EquipmentDetail.Count
		summary.Count = &count
	}

	isLiked := false
	summary.IsLiked = &isLiked

	return summary, nil
}

func toOpenAPIItemDetail(detail *domain.ItemDetail) (*openapi.ItemDetail, error) {
	if detail == nil || detail.Item == nil {
		return nil, errors.New("domain item detail is nil")
	}

	item, err := toOpenAPIItem(detail.Item)
	if err != nil {
		return nil, err
	}

	likesByUsers := make([]string, 0, len(detail.Likes))
	for _, like := range detail.Likes {
		likesByUsers = append(likesByUsers, like.UserID)
	}

	res := &openapi.ItemDetail{
		Code:        item.Code,
		Comments:    item.Comments,
		CreatedAt:   item.CreatedAt,
		DeletedAt:   item.DeletedAt,
		Description: item.Description,
		Id:          item.Id,
		ImgUrl:      item.ImgUrl,
		IsBook:      item.IsBook,
		IsTrapItem:  item.IsTrapItem,
		Name:        item.Name,
		Tags:        toOpenAPITags(detail.Tags),
		UpdatedAt:   item.UpdatedAt,
	}
	if len(likesByUsers) > 0 {
		res.LikesByUsers = &likesByUsers
	}
	if detail.Item.EquipmentDetail != nil {
		count := detail.Item.EquipmentDetail.Count
		res.Count = &count
	}

	if detail.Item.EquipmentDetail != nil {
		transactions := make([]openapi.TransactionEquipment, 0, len(detail.Transactions))
		for _, transaction := range detail.Transactions {
			transactions = append(transactions, toOpenAPITransactionEquipment(detail.Item.ID, transaction))
		}
		if err := res.FromItemDetail0(openapi.ItemDetail0{TransactionsEquipment: transactions}); err != nil {
			return nil, fmt.Errorf("failed to convert equipment transactions: %w", err)
		}
	} else {
		transactions := make([]openapi.Transaction, 0, len(detail.Transactions))
		for _, transaction := range detail.Transactions {
			transactions = append(transactions, toOpenAPITransaction(transaction))
		}
		if err := res.FromItemDetail1(openapi.ItemDetail1{Transactions: transactions}); err != nil {
			return nil, fmt.Errorf("failed to convert transactions: %w", err)
		}
	}

	return res, nil
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

func toOpenAPITags(tags []*domain.Tag) *[]openapi.Tag {
	if len(tags) == 0 {
		return nil
	}

	res := make([]openapi.Tag, 0, len(tags))
	for _, t := range tags {
		res = append(res, openapi.Tag{Name: t.Name})
	}
	return &res
}

func toOpenAPITransaction(t *domain.Transaction) openapi.Transaction {
	purpose := t.Purpose
	status := toOpenAPIBorrowingStatus(t.Status)
	userID := t.UserID

	return openapi.Transaction{
		CheckoutDate:  toOpenAPIDate(t.CheckoutDate),
		CreatedAt:     &t.CreatedAt,
		DeletedAt:     nullable.NewNullNullable[time.Time](),
		DueDate:       openapi_types.Date{Time: t.DueDate},
		Id:            &t.ID,
		Message:       t.Message,
		OwnershipId:   &t.OwnershipID,
		Purpose:       &purpose,
		ReturnMessage: t.ReturnMessage,
		ReturnDate:    toOpenAPIDate(t.ReturnDate),
		Status:        &status,
		UpdatedAt:     &t.UpdatedAt,
		UserId:        &userID,
	}
}

func toOpenAPITransactionEquipment(itemID int, t *domain.Transaction) openapi.TransactionEquipment {
	purpose := t.Purpose
	status := toOpenAPIBorrowingStatus(t.Status)
	userID := t.UserID

	return openapi.TransactionEquipment{
		CreatedAt:  &t.CreatedAt,
		DeletedAt:  nullable.NewNullNullable[time.Time](),
		DueDate:    openapi_types.Date{Time: t.DueDate},
		Id:         &t.ID,
		ItemId:     &itemID,
		Purpose:    &purpose,
		ReturnDate: toOpenAPIDate(t.ReturnDate),
		Status:     &status,
		UpdatedAt:  &t.UpdatedAt,
		UserId:     &userID,
	}
}

func toOpenAPIDate(t *time.Time) openapi_types.Date {
	if t == nil {
		return openapi_types.Date{}
	}
	return openapi_types.Date{Time: *t}
}

func toOpenAPIBorrowingStatus(status domain.BorrowingStatus) int {
	switch status {
	case domain.BorrowingStatusRequested:
		return 0
	case domain.BorrowingStatusBorrowed:
		return 1
	case domain.BorrowingStatusReturned:
		return 2
	case domain.BorrowingStatusRejected:
		return 3
	default:
		return 0
	}
}
