package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get("X-Forwarded-User")
		if userID == "" {
			return c.JSON(http.StatusUnauthorized, "X-Forwarded-User header is required")
		}

		ctx := c.Request().Context()
		ctx = WithUserID(ctx, userID)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
