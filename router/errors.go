package router

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/model"
)

func unauthorizedRequest(c echo.Context, err error) error {
	c.Logger().Infof("unauthorized request on %s: %w", c.Path(), err.Error())
	return c.String(http.StatusUnauthorized, "認証に失敗しました")
}

func forbiddenRequest(c echo.Context, err error) error {
	c.Logger().Infof("forbidden request on %s: %w", c.Path(), err.Error())
	return c.String(http.StatusForbidden, "権限がありません")
}

func invalidRequest(c echo.Context, err error) error {
	c.Logger().Infof("invalid request on %s: %w", c.Path(), err.Error())
	return c.String(http.StatusBadRequest, "リクエストデータの処理に失敗しました")
}

func internalServerError(c echo.Context, err error) error {
	c.Logger().Infof("internal server error on %s: %w", c.Path(), err.Error())
	return c.String(http.StatusInternalServerError, "予期せぬエラーが発生しました")
}

func notFoundError(c echo.Context, err error) error {
	c.Logger().Infof("not found error on %s: %w", c.Path(), err.Error())
	return c.String(http.StatusNotFound, "アイテムが見つかりません")
}

func parseModelError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, model.ErrNotFound):
		return notFoundError(c, err)
	case errors.Is(err, model.ErrUnauthorized):
		return forbiddenRequest(c, err)
	case errors.Is(err, model.ErrUpdateNotAllowed):
		return invalidRequest(c, err)
	default:
		return internalServerError(c, err)
	}
}
