package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dog/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/dog/core/dto"
	dwRepositoy "github.com/wanrun-develop/wanrun/internal/dogOwner/adapters/repository"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
)

type IDogHandler interface {
	GetAllDogs(echo.Context) ([]dto.DogListRes, error)
	GetDogByID(echo.Context, int64) (dto.DogDetailsRes, error)
	GetDogByDogOwnerID(echo.Context, int64) ([]dto.DogListRes, error)
	CreateDog(echo.Context) (dto.DogDetailsRes, error)
	DeleteDog(echo.Context, int) error
}

type dogHandler struct {
	r   repository.IDogRepository
	dwr dwRepositoy.IDogOwnerRepository
}

func NewDogHandler(r repository.IDogRepository, dwr dwRepositoy.IDogOwnerRepository) IDogHandler {
	return &dogHandler{r, dwr}
}

func (h *dogHandler) GetAllDogs(c echo.Context) ([]dto.DogListRes, error) {
	logger := log.GetLogger(c).Sugar()

	dogs, err := h.r.GetAllDogs()

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
//   - int64: 	dogのID
//
// return:
//   - dto.DogDetailsRes:	dogの詳細レスポンス
//   - error:	エラー
func (h *dogHandler) GetDogByID(c echo.Context, dogID int64) (dto.DogDetailsRes, error) {

	d, err := h.r.GetDogByID(dogID)

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

// GetDogByDogOwnerID: dogの詳細を検索して返す
//
// args:
//   - echo.Context:	コンテキスト
//   - int64: 	dogOwnerのID
//
// return:
//   - []dto.DogListRes:	dogの一覧レスポンス
//   - error:	エラー
func (h *dogHandler) GetDogByDogOwnerID(c echo.Context, dogOwnerID int64) ([]dto.DogListRes, error) {
	logger := log.GetLogger(c).Sugar()

	logger.Infof("DogOwner %d の犬の一覧検索", dogOwnerID)

	//dogownerの検索（存在チェック)
	dogOwner, err := h.dwr.GetDogOwnerById(dogOwnerID)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogOwner検索で失敗しました。", errors.NewDogServerErrorEType())
		return []dto.DogListRes{}, err
	}
	if dogOwner.IsEmpty() {
		err = errors.NewWRError(nil, "指定されたdog ownerは存在しません。", errors.NewDogClientErrorEType())
		logger.Error("不正なdog owner idでの検索")
		return []dto.DogListRes{}, err
	}

	dogs, err := h.r.GetDogByDogOwnerID(dogOwner.DogOwnerID.Int64)
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

func (h *dogHandler) CreateDog(c echo.Context) (dto.DogDetailsRes, error) {
	d, err := h.r.CreateDog()

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

func (h *dogHandler) DeleteDog(c echo.Context, dogID int) error {
	if err := h.r.DeleteDog(dogID); err != nil {
		return err
	}
	return nil
}
