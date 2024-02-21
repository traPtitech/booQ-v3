package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/traPtitech/booQ-v3/model"
	"github.com/traPtitech/booQ-v3/router"
	"github.com/traPtitech/booQ-v3/storage"
)

func main() {
	err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}
	// db.close() の必要はなさそう。参考: https://github.com/go-gorm/gorm/issues/3145

	if os.Getenv("BOOQ_ENV") == "development" {
		model.SetDBLoggerInfo()
	}

	err = model.Migrate()
	if err != nil {
		panic(err)
	}

	// Storage
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

	// Echo instance
	e := echo.New()

	// Validator
	router.SetValidator(e)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routing
	router.SetupRouting(e, router.CreateUserProvider(os.Getenv("DEBUG_USER_NAME")))

	// Start server
	e.Logger.Fatal(e.Start(":3001"))
}