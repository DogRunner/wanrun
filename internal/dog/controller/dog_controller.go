package controller

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/common"
	"github.com/wanrun-develop/wanrun/internal/dog/core/dto"
	"github.com/wanrun-develop/wanrun/internal/dog/core/handler"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IDogController interface {
	GetAllDogs(c echo.Context) error
	GetDogByID(c echo.Context) error
	GetDogByDogOwnerID(c echo.Context) error
	CreateDog(c echo.Context) error
	UpdateDog(c echo.Context) error
	DeleteDog(c echo.Context) error
}

type dogController struct {
	h handler.IDogHandler
}

func NewDogController(h handler.IDogHandler) IDogController {
	return &dogController{h}
}

func (dc *dogController) GetAllDogs(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()
	logger.Warn("dogの全検索リクエストを受け取りました。")

	resDogs, err := dc.h.GetAllDogs(c)

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
	dogID, err := strconv.ParseInt(dogIDStr, 10, 64)
	if err != nil || dogID <= 0 {
		logger.Error(err)
		err = errors.NewWRError(err, errors.M_REQUEST_PARAM_MUST_BE_NATURAL, errors.NewDogClientErrorEType())
		return err
	}
	resDog, err := dc.h.GetDogByID(c, dogID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resDog)
}

// GetDogByDogOwnerID: dogOwnerより所有している犬の一覧を取得
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - error:	エラー
func (dc *dogController) GetDogByDogOwnerID(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	dogOwnerIDStr := c.Param("dogOwnerId")
	dogOwnerID, err := strconv.ParseInt(dogOwnerIDStr, 10, 64)
	if err != nil || dogOwnerID <= 0 {
		logger.Error(err)
		err = errors.NewWRError(err, errors.M_REQUEST_PARAM_MUST_BE_NATURAL, errors.NewDogClientErrorEType())
		return err
	}

	dogs, err := dc.h.GetDogByDogOwnerID(c, dogOwnerID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dogs)
}

// CreateDog: 犬の登録
// dogIdが指定されていないこと。各フィールドのバリデーション
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - error:	エラー
func (dc *dogController) CreateDog(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	//リクエストボディをバインド
	var saveReq dto.DogSaveReq
	if err := c.Bind(&saveReq); err != nil {
		err = errors.NewWRError(err, errors.M_REQUEST_BODY_IS_INVALID, errors.NewDogrunClientErrorEType())
		logger.Error(err)
		return err
	}

	validate := validator.New()
	// カスタムバリデーションルールの登録
	_ = validate.RegisterValidation("primaryKey", common.VCreatePrimaryKey)
	_ = validate.RegisterValidation("sex", common.VSex)
	//リクエストボディのバリデーション
	if err := validate.Struct(saveReq); err != nil {
		err = errors.NewWRError(err, errors.M_REQUEST_BODY_VALIDATION_FAILED, errors.NewDogrunClientErrorEType())
		logger.Error(err)
		return err
	}

	dogId, err := dc.h.CreateDog(c, saveReq)
	if err != nil {
		return err
	}
	logger.Info("dogの作成が完了")
	return c.JSON(http.StatusOK, map[string]int64{
		"dogId": dogId,
	})
}

// UpdateDog: 犬の更新
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - error:	エラー
func (dc *dogController) UpdateDog(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	//リクエストボディをバインド
	var saveReq dto.DogSaveReq
	if err := c.Bind(&saveReq); err != nil {
		err = errors.NewWRError(err, errors.M_REQUEST_BODY_IS_INVALID, errors.NewDogrunClientErrorEType())
		logger.Error(err)
		return err
	}

	validate := validator.New()
	// カスタムバリデーションルールの登録
	_ = validate.RegisterValidation("primaryKey", common.VUpdatePrimaryKey)
	_ = validate.RegisterValidation("sex", common.VSex)
	//リクエストボディのバリデーション
	if err := validate.Struct(saveReq); err != nil {
		err = errors.NewWRError(err, errors.M_REQUEST_BODY_VALIDATION_FAILED, errors.NewDogrunClientErrorEType())
		logger.Error(err)
		return err
	}

	dogId, err := dc.h.UpdateDog(c, saveReq)
	if err != nil {
		return err
	}
	logger.Info("dogの更新が完了")
	return c.JSON(http.StatusOK, map[string]int64{
		"dogId": dogId,
	})
}

func (dc *dogController) DeleteDog(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	dogIDStr := c.Param("dogID")
	dogID, err := strconv.ParseInt(dogIDStr, 10, 64)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, errors.M_REQUEST_PARAM_MUST_BE_NATURAL, errors.NewDogClientErrorEType())
		return err
	}
	if err := dc.h.DeleteDog(c, dogID); err != nil {
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
