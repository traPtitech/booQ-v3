package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func unauthorizedRequest(c echo.Context, err error) error {
	c.Logger().Infof("unauthorized request on %s: %w", c.Path(), err.Error())
	return c.String(http.StatusUnauthorized, "認証に失敗しました")
}

func invalidRequest(c echo.Context, err error) error {
	c.Logger().Infof("invalid request on %s: %w", c.Path(), err.Error())
	return c.String(http.StatusBadRequest, "リクエストデータの処理に失敗しました")
}

func internalServerError(c echo.Context, err error) error {
	c.Logger().Infof("internal server error on %s: %w", c.Path(), err.Error())
	return c.String(http.StatusInternalServerError, "予期せぬエラーが発生しました")
}
