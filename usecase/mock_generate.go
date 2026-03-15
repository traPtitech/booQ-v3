package usecase

//go:generate mockgen -source=item.go -destination=./mock/mock_item_usecase.go -package=mock_usecase
//go:generate mockgen -source=file.go -destination=./mock/mock_file_usecase.go -package=mock_usecase
//go:generate mockgen -source=ownership.go -destination=./mock/mock_ownership_usecase.go -package=mock_usecase
