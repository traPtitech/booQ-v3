package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/booQ-v3/domain"
	"github.com/traPtitech/booQ-v3/handler/openapi"
)

func (h *handler) PostFile(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "file is required")
	}

	// ファイルを開く
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to open file: %v", err))
	}
	defer src.Close()

	// Content-Type 取得
	contentType := file.Header.Get("Content-Type")

	// UseCase 呼び出し
	uploadedFile, err := h.fu.Upload(src, contentType, file.Size)
	if err != nil {
		if errors.Is(err, domain.ErrFileTooLarge) {
			return ctx.JSON(http.StatusBadRequest, "file too large: max size is 3MB")
		}
		if errors.Is(err, domain.ErrInvalidFileType) {
			return ctx.JSON(http.StatusBadRequest, "invalid file type: only JPEG and PNG are allowed")
		}
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to upload file: %v", err))
	}

	// レスポンス
	return ctx.JSON(http.StatusCreated, openapi.File{
		Id:  uploadedFile.ID,
		Url: fmt.Sprintf("/api/files/%d", uploadedFile.ID),
	})
}

func (h *handler) GetFile(ctx echo.Context, fileId openapi.FileIdInPath) error {
	// UseCase 呼び出し
	reader, file, err := h.fu.GetFile(fileId)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to get file: %v", err))
	}
	defer reader.Close()

	// ファイルをストリームで返す
	return ctx.Stream(http.StatusOK, file.MimeType, reader)
}
