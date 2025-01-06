package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dogrunmg/core/dto"
)

type IDogrunmgHandler interface {
	DogrunmgSignUp(c echo.Context, doReq dto.DogrunmgReq) (string, error)
}

type dogrunmgHandler struct {
}

func NewDogownerHandler() IDogrunmgHandler {
	return &dogrunmgHandler{}
}

func (dmh *dogrunmgHandler) DogrunmgSignUp(c echo.Context, doReq dto.DogrunmgReq) (string, error) {
	return "", nil
}
