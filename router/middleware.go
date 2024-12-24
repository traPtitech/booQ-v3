package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var userProviderKey = "user"

// UserProvider traQに接続する用のclient
type UserProvider struct {
	AuthUser func(c echo.Context) (echo.Context, error)
}

// MiddlewareAuthUser APIにアクセスしたユーザーの情報をセットする
func (client *UserProvider) MiddlewareAuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c, err := client.AuthUser(c)
		if err != nil {
			return c.String(http.StatusUnauthorized, err.Error())
		}
		return next(c)
	}
}

func CreateUserProvider(debugUserName string) *UserProvider {
	return &UserProvider{AuthUser: func(c echo.Context) (echo.Context, error) {
		res := debugUserName
		if debugUserName == "" {
			res = c.Request().Header.Get("X-Showcase-User")
			if res == "" {
				fmt.Println(c.Request().Header)
				return c, errors.New("認証に失敗しました(Headerに必要な情報が存在しません)")
			}
		}
		c.Set(userProviderKey, res)
		return c, nil
	}}
}

func getAuthorizedUser(c echo.Context) (string, error) {
	user, ok := c.Get(userProviderKey).(string)
	if !ok {
		return "", errors.New("認証に失敗しました")
	}
	return user, nil
}
