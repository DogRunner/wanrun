package controller

import (
	"github.com/labstack/echo/v4"
	dogrunmgHandler "github.com/wanrun-develop/wanrun/internal/dogrunmg/core/handler"
)

type IDogrunmgController interface {
	DogrunmgSignUp(c echo.Context) error
}

type dogrunmgController struct {
	dm dogrunmgHandler.IDogrunmgHandler
}

func NewDogrunmgController(
	dm dogrunmgHandler.IDogrunmgHandler,
) IDogrunmgController {
	return &dogrunmgController{
		dm: dm,
	}
}

// / DogrunmgSignUp: dogrunmanagerの登録処理
//
// args:
//   - echo.Context:
//
// return:
//   - error: error情報
func (dmc *dogrunmgController) DogrunmgSignUp(c echo.Context) error {
	return nil
}
