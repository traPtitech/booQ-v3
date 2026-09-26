package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/traPtitech/booQ-v3/internal/domain"
	"github.com/traPtitech/booQ-v3/internal/handler"
	"github.com/traPtitech/booQ-v3/internal/handler/openapi"
	"github.com/traPtitech/booQ-v3/internal/middleware"
	"github.com/traPtitech/booQ-v3/internal/repository"
	"github.com/traPtitech/booQ-v3/internal/storage"
	"github.com/traPtitech/booQ-v3/internal/usecase"
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI: true,
		LogError: true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR", 
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("error", v.Error.Error()),
				)	
			}

			return nil
		},
	}))
	e.Use(echomiddleware.Recover())
	e.Use(middleware.AuthMiddleware)

	// Repository
	itemRepo := repository.NewItemRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	fileRepo := repository.NewFileRepository(db)
	ownershipRepo := repository.NewOwnershipRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	tagRepo := repository.NewTagRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	// Storage
	fileStorage := newFileStorage()

	// UseCase
	itemUseCase := usecase.NewItemUseCase(itemRepo)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, itemRepo)
	fileUseCase := usecase.NewFileUseCase(fileRepo, fileStorage)
	ownershipUseCase := usecase.NewOwnershipUseCase(ownershipRepo)
	borrowingUseCase := usecase.NewBorrowingUseCase(transactionRepo, ownershipRepo)
	tagUseCase := usecase.NewTagUseCase(tagRepo, itemRepo)
	likeUseCase := usecase.NewLikeUseCase(likeRepo, itemRepo)

	// Handler
	h := handler.NewHandlerWithTagLike(itemUseCase, commentUsecase, fileUseCase, ownershipUseCase, borrowingUseCase, tagUseCase, likeUseCase)
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
