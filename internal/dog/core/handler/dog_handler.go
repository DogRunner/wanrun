package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dog/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/dog/core/dto"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
)

type IDogHandler interface {
	GetAllDogs(echo.Context) ([]dto.DogListRes, error)
	GetDogByID(echo.Context, int) (dto.DogDetailsRes, error)
	CreateDog(echo.Context) (dto.DogDetailsRes, error)
	DeleteDog(echo.Context, int) error
}

type dogHandler struct {
	dr repository.IDogRepository
}

func NewDogHandler(dr repository.IDogRepository) IDogHandler {
	return &dogHandler{dr}
}

func (dh *dogHandler) GetAllDogs(c echo.Context) ([]dto.DogListRes, error) {
	logger := log.GetLogger(c).Sugar()

	dogs, err := dh.dr.GetAllDogs()

	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dog検索で失敗しました。", errors.NewDogServerErrorEType())
		return []dto.DogListRes{}, err
	}

	resDogs := []dto.DogListRes{}

	for _, d := range dogs {
		dr := dto.DogListRes{
			DogID:  d.DogID.Int64,
			Name:   d.Name.String,
			Weight: d.Weight.Int64,
			Sex:    d.Sex.String,
			Image:  d.Image.String,
			DogType: dto.DogTypeRes{
				DogTypeID: d.DogType.DogTypeID,
				Name:      d.DogType.Name,
			},
		}
		resDogs = append(resDogs, dr)
	}
	return resDogs, nil
}

// GetDogById: dogの詳細を検索して返す
//
// args:
//   - echo.Context:
//   - int : 	dogのID
//
// return:
//   - dto.DogDetailsRes:	dogの詳細レスポンス
//   - error:	エラー
func (dh *dogHandler) GetDogByID(c echo.Context, dogID int) (dto.DogDetailsRes, error) {

	d, err := dh.dr.GetDogByID(dogID)

	if err != nil {
		return dto.DogDetailsRes{}, err
	}

	resDog := dto.DogDetailsRes{
		DogID:      d.DogID.Int64,
		DogOwnerID: d.DogOwnerID.Int64,
		Name:       d.Name.String,
		Weight:     d.Weight.Int64,
		Sex:        d.Sex.String,
		Image:      d.Image.String,
		CreateAt:   util.ConvertToWRTime(d.CreateAt),
		UpdateAt:   util.ConvertToWRTime(d.UpdateAt),
		DogType: dto.DogTypeRes{
			DogTypeID: d.DogType.DogTypeID,
			Name:      d.DogType.Name,
		},
	}
	return resDog, nil
}

func (dh *dogHandler) CreateDog(c echo.Context) (dto.DogDetailsRes, error) {
	d, err := dh.dr.CreateDog()

	if err != nil {
		return dto.DogDetailsRes{}, err
	}

	dogRes := dto.DogDetailsRes{
		DogID:      d.DogID.Int64,
		DogOwnerID: d.DogOwnerID.Int64,
		Name:       d.Name.String,
		Weight:     d.Weight.Int64,
		Sex:        d.Sex.String,
		Image:      d.Image.String,
		CreateAt:   util.ConvertToWRTime(d.CreateAt),
		UpdateAt:   util.ConvertToWRTime(d.UpdateAt),
		DogType: dto.DogTypeRes{
			DogTypeID: d.DogType.DogTypeID,
			Name:      d.DogType.Name,
		},
	}
	return dogRes, err
}

func (dh *dogHandler) DeleteDog(c echo.Context, dogID int) error {
	if err := dh.dr.DeleteDog(dogID); err != nil {
		return err
	}
	return nil
}
