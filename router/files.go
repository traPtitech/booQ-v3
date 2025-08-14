package router

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/model"
	"github.com/traPtitech/booQ-v3/storage"
)

var uploadableMimes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// PostFile POST /files
func PostFile(c echo.Context) error {
	me, err := getAuthorizedUser(c)
	if err != nil {
		return unauthorizedRequest(c, err)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if b, ok := uploadableMimes[fileHeader.Header.Get(echo.HeaderContentType)]; !(b && ok) {
		return c.JSON(http.StatusBadRequest, errors.New("アップロードできないファイル形式です"))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return internalServerError(c, err)
	}
	defer file.Close()

	// サムネイル画像を生成
	orig, err := imaging.Decode(file, imaging.AutoOrientation(true))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New("不正なファイルです"))
	}
	newImg := image.NewRGBA(orig.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), orig, orig.Bounds().Min, draw.Over)
	b := &bytes.Buffer{}
	err = imaging.Encode(b, imaging.Fit(newImg, 360, 480, imaging.Linear), imaging.JPEG, imaging.JPEGQuality(85))
	if err != nil {
		return internalServerError(c, err)
	}

	f, err := model.CreateFile(me, b, "jpg")
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"id": f.ID, "url": fmt.Sprintf("/zpi/files/%d", f.ID)})
}

// GetFile GET /files/:id
func GetFile(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	f, err := storage.Open(fmt.Sprintf("%d.jpg", id))
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return internalServerError(c, err)
	}
	defer f.Close()

	c.Response().Header().Set("Cache-Control", "private, max-age=31536000")
	return c.Stream(http.StatusOK, "image/jpeg", f)
}
