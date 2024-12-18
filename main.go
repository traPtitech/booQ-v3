package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/traPtitech/booQ-v3/model"
	"github.com/traPtitech/booQ-v3/router"
	"github.com/traPtitech/booQ-v3/storage"
)

func main() {
	err := model.EstablishConnection()
	if err != nil {
		log.Fatal(err)
	}

	err = model.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	setStorage()

	e := echo.New()
	router.SetValidator(e)

	if os.Getenv("BOOQ_ENV") == "development" {
		e.Logger.SetLevel(log.INFO)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	router.SetupRouting(e, router.CreateUserProvider(os.Getenv("DEBUG_USER_NAME")))
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
