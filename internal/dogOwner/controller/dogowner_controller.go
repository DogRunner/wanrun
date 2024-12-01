package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authHandler "github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	"github.com/wanrun-develop/wanrun/internal/dogOwner/core/dto"
	dogOwnerHandler "github.com/wanrun-develop/wanrun/internal/dogOwner/core/handler"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IDogOwnerController interface {
	DogOwnerSignUp(c echo.Context) error
}

type dogOwnerController struct {
	doh dogOwnerHandler.IDogOwnerHandler
	ah  authHandler.IAuthHandler
}

func NewDogOwnerController(
	doh dogOwnerHandler.IDogOwnerHandler,
	ah authHandler.IAuthHandler,
) IDogOwnerController {
	return &dogOwnerController{
		doh: doh,
		ah:  ah,
	}
}

// DogOwnerSignUp: dogOwnerの登録処理
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//
// return:
//   - error: error情報
func (doc *dogOwnerController) DogOwnerSignUp(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	doReq := dto.DogOwnerReq{}

	if err := c.Bind(&doReq); err != nil {
		wrErr := errors.NewWRError(
			err,
			"入力項目に不正があります。",
			errors.NewDogOwnerClientErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	// dogOwnerのSignUp
	dogOwnerDetail, wrErr := doc.doh.DogOwnerSignUp(c, doReq)

	if wrErr != nil {
		return wrErr
	}

	// 署名済みのjwt token取得
	token, wrErr := doc.ah.GetSignedJwt(c, dogOwnerDetail)

	if wrErr != nil {
		return wrErr
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"accessToken": token,
	})
}
