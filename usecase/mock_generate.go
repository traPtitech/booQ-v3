package usecase

//go:generate mockgen -source=item.go -destination=./mock/mock_item_usecase.go -package=mock_usecase
//go:generate mockgen -source=file.go -destination=./mock/mock_file_usecase.go -package=mock_usecase
//go:generate mockgen -source=comment.go -destination=./mock/mock_comment_usecase.go -package=mock_usecase
//go:generate mockgen -source=ownership.go -destination=./mock/mock_ownership_usecase.go -package=mock_usecase
//go:generate mockgen -source=borrows.go -destination=./mock/mock_borrowing_usecase.go -package=mock_usecase
//go:generate mockgen -source=tag.go -destination=./mock/mock_tag_usecase.go -package=mock_usecase
//go:generate mockgen -source=like.go -destination=./mock/mock_like_usecase.go -package=mock_usecase
