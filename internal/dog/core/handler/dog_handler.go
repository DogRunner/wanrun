package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dog/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/dog/core/dto"
	dwRepositoy "github.com/wanrun-develop/wanrun/internal/dogOwner/adapters/repository"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
)

type IDogHandler interface {
	GetAllDogs(echo.Context) ([]dto.DogListRes, error)
	GetDogByID(echo.Context, int64) (dto.DogDetailsRes, error)
	GetDogByDogOwnerID(echo.Context, int64) ([]dto.DogListRes, error)
	CreateDog(echo.Context, dto.DogSaveReq) (int64, error)
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
		logger.Error("不正なdog owner idでの検索", err)
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

// CreateDog: 犬の登録
//
//	dogownerの存在チェック
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.DogSaveRew:	リクエスト内容
//
// return:
//   - int64:	登録したdogId
//   - error:	エラー
func (h *dogHandler) CreateDog(c echo.Context, saveReq dto.DogSaveReq) (int64, error) {
	logger := log.GetLogger(c).Sugar()

	logger.Info("create dog %w", saveReq)

	dogOwnerID := saveReq.DogOwnerID
	//dogownerの検索（存在チェック)
	dogOwner, err := h.dwr.GetDogOwnerById(dogOwnerID)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogOwner検索で失敗しました。", errors.NewDogServerErrorEType())
		return 0, err
	}
	if dogOwner.IsEmpty() {
		err = errors.NewWRError(nil, "指定されたdog ownerは存在しません。", errors.NewDogClientErrorEType())
		logger.Error("不正なdog owner idの指定", err)
		return 0, err
	}

	dog := model.Dog{
		DogOwnerID: dogOwner.DogOwnerID,
		Name:       util.NewSqlNullString(saveReq.Name),
		DogTypeID:  util.NewSqlNullInt64(saveReq.DogTypeID),
		Weight:     util.NewSqlNullInt64(saveReq.Weight),
		Sex:        util.NewSqlNullString(saveReq.Sex),
		Image:      util.NewSqlNullString(saveReq.Image),
	}
	dog, err = h.r.CreateDog(dog)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogOwnerの登録処理で失敗しました。", errors.NewDogServerErrorEType())
		return 0, err
	}

	return dog.DogID.Int64, err
}

func (h *dogHandler) DeleteDog(c echo.Context, dogID int) error {
	if err := h.r.DeleteDog(dogID); err != nil {
		return err
	}
	return nil
}
