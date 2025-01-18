package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dogrun/core/dto"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
	"gorm.io/gorm"
)

type IDogrunRepository interface {
	GetDogrunByPlaceID(echo.Context, string) (model.Dogrun, error)
	GetDogrunByID(string) (model.Dogrun, error)
	FindDogrunByIDs([]int64) ([]model.Dogrun, error)
	GetDogrunByRectanglePointerOrPlaceId(echo.Context, dto.SearchAroundRectangleCondition, []string) ([]model.Dogrun, error)
	GetDogrunByRectanglePointerAndDogrunTags(echo.Context, dto.SearchAroundRectangleCondition) ([]model.Dogrun, error)
	GetTagMst(echo.Context) ([]model.TagMst, error)
	RegistDogrunPlaceId(echo.Context, string) (int64, error)
}

type dogrunRepository struct {
	db *gorm.DB
}

func NewDogrunRepository(db *gorm.DB) IDogrunRepository {
	return &dogrunRepository{db}
}

/*
PlaceIDで、ドッグランの取得
*/
func (drr *dogrunRepository) GetDogrunByPlaceID(c echo.Context, placeID string) (model.Dogrun, error) {
	logger := log.GetLogger(c).Sugar()
	dogrun := model.Dogrun{}
	if err := drr.db.Preload("DogrunTags").
		Preload("RegularBusinessHours").
		Preload("SpecialBusinessHours").
		Where("place_id = ?", placeID).
		Find(&dogrun).Error; err != nil {
		logger.Error(err)
		return model.Dogrun{}, errors.NewWRError(err, "DBからのデータ取得に失敗", errors.NewDogrunServerErrorEType())
	}
	return dogrun, nil
}

/*
DogrunIDで、ドッグランの取得
*/
func (drr *dogrunRepository) GetDogrunByID(id string) (model.Dogrun, error) {
	dogrun := model.Dogrun{}
	if err := drr.db.Where("dogrun_id = ?", id).Find(&dogrun).Error; err != nil {
		return dogrun, err
	}
	return dogrun, nil
}

// FindDogrunByIDs: 複数IDのドッグラン検索
//
// args:
//   - echo.Context:	コンテキスト
//   - []int64: dogrunIDs
//
// return:
//   - []model.Dogrun:	検索結果
//   - error:	エラー
func (drr *dogrunRepository) FindDogrunByIDs(ids []int64) ([]model.Dogrun, error) {
	dogruns := []model.Dogrun{}
	if err := drr.db.Where("dogrun_id in ?", ids).Find(&dogruns).Error; err != nil {
		return dogruns, err
	}
	return dogruns, nil
}

// GetDogrunByRectanglePointerOrPlaceId: 条件の範囲内 または 指定のPlaceIDのdogrunを取得
//
// args:
//   - echo.Context:	コンテキスト
//   - to.SearchAroundRectangleCondition:	条件
//   - []string:	placeIDs
//
// return:
//   - []model.Dogrun:	ドッグランの検索結果
//   - error:	エラー
func (drr *dogrunRepository) GetDogrunByRectanglePointerOrPlaceId(c echo.Context, condition dto.SearchAroundRectangleCondition, placeIDs []string) ([]model.Dogrun, error) {
	logger := log.GetLogger(c).Sugar()
	dogruns := []model.Dogrun{}
	if err := drr.db.Preload("DogrunTags").
		Preload("RegularBusinessHours").
		Preload("SpecialBusinessHours").
		Where("(longitude BETWEEN ? AND ?) AND (latitude BETWEEN ? AND ?)",
			condition.Target.Southwest.Longitude, condition.Target.Northeast.Longitude,
			condition.Target.Southwest.Latitude, condition.Target.Northeast.Latitude).
		Or("place_id IN ?", placeIDs).
		Find(&dogruns).Error; err != nil {
		logger.Error(err)
		return nil, errors.NewWRError(err, "DBからのデータ取得に失敗", errors.NewDogrunServerErrorEType())
	}
	return dogruns, nil
}

// GetDogrunByRectanglePointerAndDogrunTags: 条件の範囲内 かつ ドッグランタグのdogrunを取得
//
// args:
//   - echo.Context:	コンテキスト
//   - to.SearchAroundRectangleCondition:	条件
//
// return:
//   - []model.Dogrun:	ドッグランの検索結果
//   - error:	エラー
func (drr *dogrunRepository) GetDogrunByRectanglePointerAndDogrunTags(c echo.Context, condition dto.SearchAroundRectangleCondition) ([]model.Dogrun, error) {
	logger := log.GetLogger(c).Sugar()
	dogruns := []model.Dogrun{}
	if err := drr.db.Joins("LEFT OUTER JOIN dogrun_tags on dogruns.dogrun_id = dogrun_tags.dogrun_id").
		Where("(longitude BETWEEN ? AND ?) AND (latitude BETWEEN ? AND ?)",
			condition.Target.Southwest.Longitude, condition.Target.Northeast.Longitude,
			condition.Target.Southwest.Latitude, condition.Target.Northeast.Latitude).
		Where("dogrun_tags.tag_id IN ?", condition.IncludeDogrunTags).
		Group("dogruns.dogrun_id"). // dogruns の重複を排除
		Preload("DogrunTags").
		Preload("RegularBusinessHours").
		Preload("SpecialBusinessHours").
		Find(&dogruns).Error; err != nil {
		logger.Error(err)
		return nil, errors.NewWRError(err, "DBからのデータ取得に失敗", errors.NewDogrunServerErrorEType())
	}
	return dogruns, nil
}

// GetDogrunTagMst: tag_mstの全件select
//
// args:
//   - echo.context:	コンテキスト
//
// return:
//   - []model.TagMst:	マスター情報
//   - error:	エラー
func (drr *dogrunRepository) GetTagMst(c echo.Context) ([]model.TagMst, error) {
	logger := log.GetLogger(c).Sugar()

	tagMst := []model.TagMst{}
	if err := drr.db.Find(&tagMst).Error; err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "tag_mstのselectで失敗しました。", errors.NewDogServerErrorEType())
		return []model.TagMst{}, err
	}
	return tagMst, nil
}

// RegistDogrunPlaceId: placeIdをDBへ保存する
//
// args:
//   - echo.Context: c
//   - string: placeId	DBへ保存するplaceId
//
// return:
//   - int:	dogrunsテーブルのPK
//   - error:	エラー
func (drr *dogrunRepository) RegistDogrunPlaceId(c echo.Context, placeId string) (int64, error) {
	logger := log.GetLogger(c).Sugar()
	dogrun := model.Dogrun{PlaceId: util.NewSqlNullString(placeId)}

	if err := drr.db.Create(&dogrun).Error; err != nil {
		logger.Error(err)
		err := errors.NewWRError(err, "placeIdのDB保存に失敗", errors.NewDogrunServerErrorEType())
		return 0, err
	}

	//主キー返す
	return dogrun.DogrunID.Int64, nil
}
