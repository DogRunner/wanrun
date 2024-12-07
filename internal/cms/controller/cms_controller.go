package controller

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/cms/core/dto"
	"github.com/wanrun-develop/wanrun/internal/cms/core/handler"
	"github.com/wanrun-develop/wanrun/internal/wrcontext"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type ICmsController interface {
	UploadFile(c echo.Context) error
}

type cmsController struct {
	ch handler.ICmsHandler
}

func NewCmsController(ch handler.ICmsHandler) ICmsController {
	return &cmsController{ch}
}

func (cc *cmsController) UploadFile(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	// claims情報の取得
	claims, wrErr := wrcontext.GetVerifiedClaims(c)

	if wrErr != nil {
		return wrErr
	}

	// dogOwnerIDの取得
	dogOwnerID, wrErr := claims.GetDogOwnerIDAsInt64(c)
	if wrErr != nil {
		return wrErr
	}

	// フォームからファイルの取得
	file, err := c.FormFile("file") // "file"はフロントエンドのフォームデータのキー

	if err != nil {
		wrErr := errors.NewWRError(
			err,
			"ファイルデータに不正があります。",
			errors.NewDogOwnerClientErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	// ファイルの内容を開く
	src, err := file.Open()

	if err != nil {
		wrErr := errors.NewWRError(
			err,
			"ファイルデータに不正があります。",
			errors.NewDogOwnerClientErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}
	defer src.Close()

	// ファイル名を取得
	fileName := file.Filename

	// 拡張子を取得
	ext := strings.TrimPrefix(filepath.Ext(fileName), ".")

	// 拡張子を除いたファイル名を取得
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	fuq := dto.FileUploadReq{
		FileName:   baseName,
		Extension:  ext,
		Src:        src,
		DogOwnerID: dogOwnerID,
	}

	// FileUploadのハンドラー
	if wrErr := cc.ch.HandleFileUpload(c, fuq); wrErr != nil {
		return wrErr
	}

	return c.JSON(http.StatusOK, map[string]any{})
}
