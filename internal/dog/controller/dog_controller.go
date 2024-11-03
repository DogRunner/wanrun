package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dog/core/handler"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IDogController interface {
	GetAllDogs(c echo.Context) error
	GetDogByID(c echo.Context) error
	CreateDog(c echo.Context) error
	DeleteDog(c echo.Context) error
}

type dogController struct {
	dh handler.IDogHandler
}

func NewDogController(dh handler.IDogHandler) IDogController {
	return &dogController{dh}
}

func (dc *dogController) GetAllDogs(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()
	logger.Warn("dogの全検索リクエストを受け取りました。")

	resDogs, err := dc.dh.GetAllDogs(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resDogs)
}

// GetDogById: 犬の詳細を取得
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - error:	エラー
func (dc *dogController) GetDogByID(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	dogIDStr := c.Param("dogID")
	dogID, err := strconv.Atoi(dogIDStr)
	if err != nil || dogID > 0 {
		logger.Error(err)
		err = errors.NewWRError(err, "このリクエストパラメーターには整数のみ指定可能です。", errors.NewDogClientErrorEType())
		return err
	}
	resDog, err := dc.dh.GetDogByID(c, dogID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resDog)
}

func (dc *dogController) CreateDog(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	resDog, err := dc.dh.CreateDog(c)
	if err != nil {
		return err
	}
	logger.Info("dogの作成が完了")
	return c.JSON(http.StatusOK, resDog)
}

func (dc *dogController) DeleteDog(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	dogIDStr := c.Param("dogID")
	dogID, err := strconv.Atoi(dogIDStr)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "このリクエストパラメーターには整数のみ指定可能です。", errors.NewDogClientErrorEType())
		return err
	}
	if err := dc.dh.DeleteDog(c, dogID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// func (dc *dogController) UpdateDog(c echo.Context) error {
// 	dogIDStr := c.Param("dogID")
// 	dogID, err := strconv.Atoi(dogIDStr)
// 	if err != nil {
// 		log.Error(err)
// 		return c.JSON(http.StatusBadRequest, errors.ErrorResponse{
// 			Code:    http.StatusBadRequest,
// 			Message: "Invalid dog ID format",
// 		})
// 	}
// }
