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
		panic(err)
	}

	err = model.Migrate()
	if err != nil {
		panic(err)
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
	if os.Getenv("OS_CONTAINER") != "" {
		// Swiftオブジェクトストレージ
		err := storage.SetSwiftStorage(
			os.Getenv("OS_CONTAINER"),
			os.Getenv("OS_USERNAME"),
			os.Getenv("OS_PASSWORD"),
			os.Getenv("OS_TENANT_NAME"),
			os.Getenv("OS_TENANT_ID"),
			os.Getenv("OS_AUTH_URL"),
		)
		if err != nil {
			panic(err)
		}
	} else {
		// ローカルストレージ
		dir := os.Getenv("UPLOAD_DIR")
		if dir == "" {
			dir = "./uploads"
		}
		err := storage.SetLocalStorage(dir)
		if err != nil {
			panic(err)
		}
	}
}
