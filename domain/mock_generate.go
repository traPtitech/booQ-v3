package domain

//go:generate mockgen -source=item.go -destination=./mock/mock_item_repository.go -package=mock_domain
//go:generate mockgen -source=ownership.go -destination=./mock/mock_ownership_repository.go -package=mock_domain
//go:generate mockgen -source=borrows.go -destination=./mock/mock_transaction_repository.go -package=mock_domain
//go:generate mockgen -source=tag.go -destination=./mock/mock_tag_repository.go -package=mock_domain
//go:generate mockgen -source=like.go -destination=./mock/mock_like_repository.go -package=mock_domain
