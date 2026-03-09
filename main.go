package main

import (
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/middleware"
	"github.com/traPtitech/booQ-v3/repository"
	"github.com/traPtitech/booQ-v3/storage"
	"github.com/traPtitech/booQ-v3/usecase"
)

func main() {
	db, err := repository.EstablishConnection()
	if err != nil {
		log.Fatal(err)
	}

	err = repository.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	if os.Getenv("BOOQ_ENV") == "development" {
		repository.SetLoggerInfo(db)
		e.Logger.SetLevel(log.INFO)
	}

	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(middleware.AuthMiddleware)

	// Repository
	itemRepo := repository.NewItemRepository(db)
	fileRepo := repository.NewFileRepository(db)
	ownershipRepo := repository.NewOwnershipRepository(db)

	// Storage
	fileStorage := newFileStorage()

	// UseCase
	itemUseCase := usecase.NewItemUseCase(itemRepo)
	fileUseCase := usecase.NewFileUseCase(fileRepo, fileStorage)
	ownershipUseCase := usecase.NewOwnershipUseCase(ownershipRepo)

	// Handler
	h := handler.NewHandler(itemUseCase, fileUseCase, ownershipUseCase)
	openapi.RegisterHandlers(e, h)

	e.Logger.Fatal(e.Start(":3001"))
}

func newFileStorage() domain.FileStorage {
	if os.Getenv("S3_BUCKET") != "" {
		// S3
		s, err := storage.NewS3Storage(
			os.Getenv("S3_BUCKET"),
			os.Getenv("S3_REGION"),
			os.Getenv("S3_ENDPOINT"),
			os.Getenv("S3_ACCESS_KEY"),
			os.Getenv("S3_SECRET_KEY"),
		)
		if err != nil {
			log.Fatal(err)
		}
		return s
	}

	// ローカルストレージ
	dir := os.Getenv("UPLOAD_DIR")
	if dir == "" {
		dir = "./uploads"
	}
	s, err := storage.NewLocalStorage(dir)
	if err != nil {
		log.Fatal(err)
	}
	return s
}
