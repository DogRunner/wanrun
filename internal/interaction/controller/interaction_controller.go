package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/interaction/core/dto"
	"github.com/wanrun-develop/wanrun/internal/interaction/core/handler"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IBookmarkController interface {
	AddBookmark(echo.Context) error
}

type bookmarkController struct {
	h handler.IBookmarkHandler
}

func NewBookmarkController(bh handler.IBookmarkHandler) IBookmarkController {
	return &bookmarkController{bh}
}

// AddBookmark: ブックマークの追加
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - error:	エラー
func (bc *bookmarkController) AddBookmark(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	var reqBody dto.AddBookmark
	if err := c.Bind(&reqBody); err != nil {
		err = errors.NewWRError(err, "ブックマーク登録リクエストが不正です", errors.NewInteractionClientErrorEType())
		logger.Error(err)
		return err
	}
	// バリデータのインスタンス作成
	validate := validator.New()

	//リクエストボディのバリデーション
	if err := validate.Struct(reqBody); err != nil {
		err = errors.NewWRError(err, "リクエストがバリデーションに違反しています", errors.NewInteractionClientErrorEType())
		logger.Error(err)
		return err
	}

	//本処理
	bookmarkId, err := bc.h.AddBookmark(c, reqBody)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"dogrunBookmarkId": bookmarkId,
	})
}
