package router

import (
	"net/http"

	"github.com/labstack/echo"
)

// アップロードを許可するMIMEタイプ
var uploadableMimes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// PostFile POST /files
func PostFile(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented");
}

// GetFile GET /files/:id
func GetFile(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented");
}
