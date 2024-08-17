package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dogrun/core/handler"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IDogrunController interface {
	GetDogrunDetail(echo.Context) error
}

type DogrunController struct {
	h handler.IDogrunHandler
}

func NewDogrunController(h handler.IDogrunHandler) IDogrunController {
	return &DogrunController{h}
}

func (dc *DogrunController) GetDogrunDetail(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	placeId := c.Param("placeId")
	logger.Info("リクエストplace id :", placeId)
	log.GetLogger(c).Info("リクエストplace id")

	dogrun, err := dc.h.GetDogrunDetail(c, placeId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve dog run information",
		})
	}

	return c.JSON(http.StatusOK, dogrun)
}