package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/common"
	"github.com/wanrun-develop/wanrun/internal/interaction/core/dto"
	"github.com/wanrun-develop/wanrun/internal/interaction/core/handler"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IInteractionController interface {
	AddBookmark(echo.Context) error
	DeleteBookmarks(echo.Context) error
	CheckinDogrun(echo.Context) error
	CheckoutDogrun(echo.Context) error
}

type interactionController struct {
	bh handler.IBookmarkHandler
	ch handler.ICheckInOutHandler
}

func NewInteractionController(bh handler.IBookmarkHandler, ch handler.ICheckInOutHandler) IInteractionController {
	return &interactionController{bh, ch}
}

// AddBookmark: ブックマークの追加
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - error:	エラー
func (ic *interactionController) AddBookmark(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	reqBody := dto.BookmarkAddReq{}
	if err := c.Bind(&reqBody); err != nil {
		err = errors.NewWRError(err, "ブックマーク登録リクエストが不正です", errors.NewInteractionClientErrorEType())
		logger.Error(err)
		return err
	}
	// バリデータのインスタンス作成
	validate := validator.New()
	_ = validate.RegisterValidation("notEmpty", common.VNotEmpty)

	//リクエストボディのバリデーション
	if err := validate.Struct(reqBody); err != nil {
		err = errors.NewWRError(err, "リクエストがバリデーションに違反しています", errors.NewInteractionClientErrorEType())
		logger.Error(err)
		return err
	}

	//本処理
	bookmarkId, err := ic.bh.AddBookmark(c, reqBody)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string][]int64{
		"dogrunBookmarkId": bookmarkId,
	})
}

// DeleteBookmarks: ブックマーク削除
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - error	:	エラー
func (ic *interactionController) DeleteBookmarks(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	reqBody := dto.BookmarkDeleteReq{}
	if err := c.Bind(&reqBody); err != nil {
		err = errors.NewWRError(err, "ブックマーク登録リクエストが不正です", errors.NewInteractionClientErrorEType())
		logger.Error(err)
		return err
	}
	// バリデータのインスタンス作成
	validate := validator.New()
	_ = validate.RegisterValidation("notEmpty", common.VNotEmpty)
	//リクエストボディのバリデーション
	if err := validate.Struct(reqBody); err != nil {
		err = errors.NewWRError(err, "リクエストがバリデーションに違反しています", errors.NewInteractionClientErrorEType())
		logger.Error(err)
		return err
	}

	if err := ic.bh.DeleteBookmark(c, reqBody); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)

}

// CheckinDogrun: ドッグランへのチェックイン（入場記録）
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
// error:	　エラー
func (ic *interactionController) CheckinDogrun(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	reqBody := dto.CheckinReq{}
	if err := c.Bind(&reqBody); err != nil {
		err = errors.NewWRError(err, "チェックインリクエストが不正です", errors.NewInteractionClientErrorEType())
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

	err := ic.ch.CheckinDogrun(c, reqBody)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

// CheckoutDogrun: ドッグランへのチェックアウト（退場記録）
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
// error:	　エラー
func (ic *interactionController) CheckoutDogrun(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	reqBody := dto.CheckoutReq{}
	if err := c.Bind(&reqBody); err != nil {
		err = errors.NewWRError(err, "チェックアウトリクエストが不正です", errors.NewInteractionClientErrorEType())
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

	err := ic.ch.CheckoutDogrun(c, reqBody)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
