package router

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

// SetupRouting APIのルーティングを行います
func SetupRouting(e *echo.Echo, client *UserProvider) {
	e.GET("/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := e.Group("/api", client.MiddlewareAuthUser)
	{
		apiItems := api.Group("/items")
		{
			apiItems.GET("", GetItems)
			apiItems.POST("", PostItems)
			apiItems.GET("/:id", GetItem)
			apiItems.PUT("/:id", PutItem)
			apiItems.DELETE("/:id", DeleteItem)

			apiItems.POST("/:id/owners", PostOwners)
			apiItems.PATCH("/:id/owners/:ownershipid", PatchOwners)
			apiItems.DELETE("/:id/owners/:ownershipid", DeleteOwners)
			apiItems.POST("/:id/comments", PostComments)
			apiItems.POST("/:id/likes", PostLikes)
			apiItems.DELETE("/:id/likes", DeleteLikes)

			apiBorrowingEquipment := apiItems.Group("/:id/borrowing")
			{
				apiBorrowingEquipment.POST("", PostBorrowingEquipment)
				apiBorrowingEquipment.POST("/return", PostBorrowingEquipmentReturn)
			}

			apiOwnersBorrowing := apiItems.Group("/:id/owners/:ownershipid/borrowing")
			{
				apiOwnersBorrowing.POST("", PostBorrowings)
				apiOwnersBorrowing.GET("/:borrowingid", GetBorrowingsId)
				apiOwnersBorrowing.POST("/:borrowingid/reply", PostBorrowingsReply)
				apiOwnersBorrowing.POST("/:borrowingid/return", PostBorrowingsReturn)
			}
		}

		apiFiles := api.Group("/files")
		{
			apiFiles.POST("", PostFile, middleware.BodyLimit("3MB"))
		}

	}
	e.GET("/api/files/:id", GetFile)
}
