package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/traPtitech/booQ-v3/handler"
	"github.com/traPtitech/booQ-v3/handler/openapi"
	"github.com/traPtitech/booQ-v3/repository"
	"github.com/traPtitech/booQ-v3/usecase"

	"github.com/traPtitech/booQ-v3/storage"
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

	setStorage()

	e := echo.New()

	if os.Getenv("BOOQ_ENV") == "development" {
		repository.SetLoggerInfo(db)
		e.Logger.SetLevel(log.INFO)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Repository
	itemRepo := repository.NewItemRepository(db)
	fileRepo := repository.NewFileRepository(db)

	// Storage
	fileStorage := storage.NewFileStorage()

	// UseCase
	itemUseCase := usecase.NewItemUseCase(itemRepo)
	fileUseCase := usecase.NewFileUseCase(fileRepo, fileStorage)

	// Handler
	h := handler.NewHandler(itemUseCase, fileUseCase)
	openapi.RegisterHandlers(e, h)

	e.Logger.Fatal(e.Start(":3001"))
}

func setStorage() {
	if os.Getenv("S3_BUCKET") != "" {
		// S3
		err := storage.SetS3Storage(
			os.Getenv("S3_BUCKET"),
			os.Getenv("S3_REGION"),
			os.Getenv("S3_ENDPOINT"),
			os.Getenv("S3_ACCESS_KEY"),
			os.Getenv("S3_SECRET_KEY"),
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// ローカルストレージ
		dir := os.Getenv("UPLOAD_DIR")
		if dir == "" {
			dir = "./uploads"
		}
		err := storage.SetLocalStorage(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
}
