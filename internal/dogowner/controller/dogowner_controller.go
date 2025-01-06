package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	authHandler "github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	doDTO "github.com/wanrun-develop/wanrun/internal/dogowner/core/dto"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IDogownerController interface {
	DogownerSignUp(c echo.Context) error
}

type dogownerController struct {
	doh dogownerHandler.IDogownerHandler
	ah  authHandler.IAuthHandler
}

func NewDogownerController(
	doh dogownerHandler.IDogownerHandler,
	ah authHandler.IAuthHandler,
) IDogownerController {
	return &dogownerController{
		doh: doh,
		ah:  ah,
	}
}

// DogownerSignUp: dogownerの登録処理
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//
// return:
//   - error: error情報
func (doc *dogownerController) DogownerSignUp(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	doReq := doDTO.DogownerReq{}

	if err := c.Bind(&doReq); err != nil {
		wrErr := errors.NewWRError(
			err,
			"入力項目に不正があります。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	// バリデータのインスタンス作成
	validate := validator.New()

	//リクエストボディのバリデーション
	if err := validate.Struct(&doReq); err != nil {
		err = errors.NewWRError(
			err,
			"必須の項目に不正があります。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(err)
		return err
	}

	// dogownerのSignUp
	token, wrErr := doc.doh.DogownerSignUp(c, doReq)

	if wrErr != nil {
		return wrErr
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"accessToken": token,
	})
}
