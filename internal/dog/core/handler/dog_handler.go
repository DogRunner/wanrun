package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dog/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/dog/core/dto"
	dwRepository "github.com/wanrun-develop/wanrun/internal/dogowner/adapters/repository"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
)

type IDogHandler interface {
	GetAllDogs(echo.Context) ([]dto.DogListRes, error)
	GetDogByID(echo.Context, int64) (dto.DogDetailsRes, error)
	GetDogByDogownerID(echo.Context, int64) ([]dto.DogListRes, error)
	GetDogTypeMst(c echo.Context) ([]dto.DogTypeMstRes, error)
	CreateDog(echo.Context, dto.DogSaveReq) (int64, error)
	UpdateDog(echo.Context, dto.DogSaveReq) (int64, error)
	DeleteDog(echo.Context, int64) error
}

type dogHandler struct {
	r   repository.IDogRepository
	dwr dwRepository.IDogownerRepository
}

func NewDogHandler(r repository.IDogRepository, dwr dwRepository.IDogownerRepository) IDogHandler {
	return &dogHandler{r, dwr}
}

func (h *dogHandler) GetAllDogs(c echo.Context) ([]dto.DogListRes, error) {
	logger := log.GetLogger(c).Sugar()

	dogs, err := h.r.GetAllDogs(c)

	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dog検索で失敗しました。", errors.NewDogServerErrorEType())
		return []dto.DogListRes{}, err
	}

	resDogs := []dto.DogListRes{}

	for _, d := range dogs {
		dr := dto.DogListRes{
			DogID:     d.DogID.Int64,
			Name:      d.Name.String,
			Weight:    d.Weight.Int64,
			Sex:       d.Sex.String,
			Image:     d.Image.String,
			DogTypeId: []int64{d.DogTypeID.Int64},
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

	d, err := h.r.GetDogByID(c, dogID)

	if err != nil {
		return dto.DogDetailsRes{}, err
	}

	resDog := dto.DogDetailsRes{
		DogID:      d.DogID.Int64,
		DogownerID: d.DogownerID.Int64,
		Name:       d.Name.String,
		Weight:     d.Weight.Int64,
		Sex:        d.Sex.String,
		Image:      d.Image.String,
		DogTypeId:  []int64{d.DogTypeID.Int64},
		CreateAt:   util.ConvertToWRTime(d.CreateAt),
		UpdateAt:   util.ConvertToWRTime(d.UpdateAt),
	}
	return resDog, nil
}

// GetDogByDogownerID: dogの詳細を検索して返す
//
// args:
//   - echo.Context:	コンテキスト
//   - int64: 	dogownerのID
//
// return:
//   - []dto.DogListRes:	dogの一覧レスポンス
//   - error:	エラー
func (h *dogHandler) GetDogByDogownerID(c echo.Context, dogownerID int64) ([]dto.DogListRes, error) {
	logger := log.GetLogger(c).Sugar()

	logger.Infof("Dogowner %d の犬の一覧検索", dogownerID)

	//dogownerの検索（存在チェック)
	dogowner, err := h.dwr.GetDogownerById(dogownerID)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogowner検索で失敗しました。", errors.NewDogServerErrorEType())
		return []dto.DogListRes{}, err
	}
	if dogowner.IsEmpty() {
		err = errors.NewWRError(nil, "指定されたdog ownerは存在しません。", errors.NewDogClientErrorEType())
		logger.Error("不正なdog owner idでの検索", err)
		return []dto.DogListRes{}, err
	}

	dogs, err := h.r.GetDogByDogownerID(c, dogowner.DogownerID.Int64)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dog検索で失敗しました。", errors.NewDogServerErrorEType())
		return []dto.DogListRes{}, err
	}

	resDogs := []dto.DogListRes{}

	for _, d := range dogs {
		dr := dto.DogListRes{
			DogID:     d.DogID.Int64,
			Name:      d.Name.String,
			Weight:    d.Weight.Int64,
			Sex:       d.Sex.String,
			Image:     d.Image.String,
			DogTypeId: []int64{d.DogTypeID.Int64},
		}
		resDogs = append(resDogs, dr)
	}

	return resDogs, nil
}

// GetDogTypeMst: DogTypeマスター情報の取得
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - []dto.DogTypeMstRes:	マスター情報
//   - error:	エラー
func (h *dogHandler) GetDogTypeMst(c echo.Context) ([]dto.DogTypeMstRes, error) {
	dogTypeMst, err := h.r.GetDogTypeMst(c)
	if err != nil {
		return []dto.DogTypeMstRes{}, err
	}
	mstRes := []dto.DogTypeMstRes{}

	for _, m := range dogTypeMst {
		mst := dto.DogTypeMstRes{
			DogTypeID: m.DogTypeID,
			Name:      m.Name,
		}
		mstRes = append(mstRes, mst)
	}

	return mstRes, nil
}

// CreateDog: 犬の登録
//
//	dogownerの存在チェック
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.DogSaveReq:	リクエスト内容
//
// return:
//   - int64:	登録したdogId
//   - error:	エラー
func (h *dogHandler) CreateDog(c echo.Context, saveReq dto.DogSaveReq) (int64, error) {
	logger := log.GetLogger(c).Sugar()

	logger.Info("create dog %v", saveReq)

	dogownerID := saveReq.DogownerID
	//dogownerの検索（存在チェック)
	if err := h.isExistsDogowner(c, dogownerID); err != nil {
		return 0, err
	}

	dog := model.Dog{
		DogownerID: util.NewSqlNullInt64(dogownerID),
		Name:       util.NewSqlNullString(saveReq.Name),
		DogTypeID:  util.NewSqlNullInt64(saveReq.DogTypeID),
		Weight:     util.NewSqlNullInt64(saveReq.Weight),
		Sex:        util.NewSqlNullString(saveReq.Sex),
		Image:      util.NewSqlNullString(saveReq.Image),
	}

	dog, err := h.r.CreateDog(c, dog)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogの登録処理で失敗しました。", errors.NewDogServerErrorEType())
		return 0, err
	}

	return dog.DogID.Int64, err
}

// UpdateDog: dogの更新
//
//	dogの存在チェック
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.DogSaveReq:	リクエスト内容
//
// return:
//   - int64:	更新したdogID
//   - error:	エラー
func (h *dogHandler) UpdateDog(c echo.Context, saveReq dto.DogSaveReq) (int64, error) {
	logger := log.GetLogger(c).Sugar()

	logger.Info("update dog %v", saveReq)

	dogID := saveReq.DogID
	// dogの存在チェック
	var dog model.Dog
	var err error
	if dog, err = h.isExistsDog(c, dogID); err != nil {
		return 0, err
	}

	//dogownerが変わっていれば存在チェック
	if saveReq.DogownerID != dog.DogownerID.Int64 {
		dogownerID := saveReq.DogownerID
		if err = h.isExistsDogowner(c, dogownerID); err != nil {
			return 0, err
		}
	}

	//更新値をつめる
	dog.DogownerID = util.NewSqlNullInt64(saveReq.DogownerID)
	dog.Name = util.NewSqlNullString(saveReq.Name)
	dog.DogTypeID = util.NewSqlNullInt64(saveReq.DogTypeID)
	dog.Weight = util.NewSqlNullInt64(saveReq.Weight)
	dog.Sex = util.NewSqlNullString(saveReq.Sex)
	dog.Image = util.NewSqlNullString(saveReq.Image)
	//更新
	dog, err = h.r.UpdateDog(c, dog)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogの更新処理で失敗しました。", errors.NewDogServerErrorEType())
		return 0, err
	}

	return dog.DogID.Int64, err
}

func (h *dogHandler) DeleteDog(c echo.Context, dogID int64) error {
	if _, err := h.isExistsDog(c, dogID); err != nil {
		return err
	}
	if err := h.r.DeleteDog(c, dogID); err != nil {
		return err
	}
	return nil
}

// isExistsDog: dogの存在チェック
//
// args:
//   - echo.Context:	コンテキスト
//   - int64:	チェック対象のdogID
//
// return:
//   - error:	エラー
func (h *dogHandler) isExistsDog(c echo.Context, dogID int64) (model.Dog, error) {
	logger := log.GetLogger(c).Sugar()

	dog, err := h.r.GetDogByID(c, dogID)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dog検索で失敗しました。", errors.NewDogServerErrorEType())
		return model.Dog{}, err
	}
	if dog.IsEmpty() {
		err = errors.NewWRError(nil, "指定されたdogは存在しません。", errors.NewDogClientErrorEType())
		logger.Error("不正なdog idの指定", err)
		return model.Dog{}, err
	}
	return dog, nil
}

// isExistsDogowner: dogownerの存在チェック
//
// args:
//   - echo.Context:	コンテキスト
//   - int64:	チェック対象のdogownerId
//
// return:
//   - error:	エラー
func (h *dogHandler) isExistsDogowner(c echo.Context, dogownerID int64) error {
	logger := log.GetLogger(c).Sugar()
	//検索
	dogowner, err := h.dwr.GetDogownerById(dogownerID)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogowner検索で失敗しました。", errors.NewDogServerErrorEType())
		return err
	}
	if dogowner.IsEmpty() {
		err = errors.NewWRError(nil, "指定されたdog ownerは存在しません。", errors.NewDogClientErrorEType())
		logger.Error("不正なdog owner idの指定", err)
		return err
	}
	return nil
}
