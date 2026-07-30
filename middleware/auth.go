package middleware

import "github.com/labstack/echo/v4"

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ctx = WithUserID(ctx, "sample-user") // TODO: get user ID

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
