package repository

import (
	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IDogRepository interface {
	GetAllDogs(echo.Context) ([]model.Dog, error)
	GetDogByID(echo.Context, int64) (model.Dog, error)
	GetDogByDogOwnerID(echo.Context, int64) ([]model.Dog, error)
	GetDogTypeMst(echo.Context) ([]model.DogTypeMst, error)
	CreateDog(echo.Context, model.Dog) (model.Dog, error)
	UpdateDog(echo.Context, model.Dog) (model.Dog, error)
	DeleteDog(echo.Context, int64) error
}

type dogRepository struct {
	db *gorm.DB
}

func NewDogRepository(db *gorm.DB) IDogRepository {
	return &dogRepository{db}
}

func (dr *dogRepository) GetAllDogs(c echo.Context) ([]model.Dog, error) {
	logger := log.GetLogger(c).Sugar()

	dogs := []model.Dog{}
	if err := dr.db.Find(&dogs).Error; err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogのselectで失敗しました。", errors.NewDogServerErrorEType())
		return []model.Dog{}, err
	}
	return dogs, nil
}

// GetDogByID: DBへDogIDでdogsのセレクト。dogTypeもロードする
//
// args:
//   - int64:	dogId
//
// return:
//   - model.Dog:	dogデータ
//   - error:	エラー
func (dr *dogRepository) GetDogByID(c echo.Context, dogID int64) (model.Dog, error) {
	logger := log.GetLogger(c).Sugar()

	dog := model.Dog{}
	if err := dr.db.Where("dog_id=?", dogID).Find(&dog).Error; err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogのselectで失敗しました。", errors.NewDogServerErrorEType())
		return model.Dog{}, err
	}
	return dog, nil
}

// GetDogByID: DBへDogOwnerIDでdogsのセレクト。dogTypeもロードする
//
// args:
//   - int:	dogId
//
// return:
//   - []model.Dog:	dogデータ
//   - error:	エラー
func (dr *dogRepository) GetDogByDogOwnerID(c echo.Context, dogOwnerID int64) ([]model.Dog, error) {
	logger := log.GetLogger(c).Sugar()

	dogs := []model.Dog{}
	if err := dr.db.Where("dog_owner_id=?", dogOwnerID).Find(&dogs).Error; err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogのselectで失敗しました。", errors.NewDogServerErrorEType())
		return []model.Dog{}, err
	}
	return dogs, nil
}

// GetDogTypeMst: dog_type_mstからマスターデータの全権select
//
// args:
//   - echo.Context:
//     -:
//
// return:
//   - []model.DogTypeMst:
//   - error:
func (dr *dogRepository) GetDogTypeMst(c echo.Context) ([]model.DogTypeMst, error) {
	logger := log.GetLogger(c).Sugar()

	dogTypeMst := []model.DogTypeMst{}
	if err := dr.db.Find(&dogTypeMst).Error; err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dog_type_mstのselectで失敗しました。", errors.NewDogServerErrorEType())
		return []model.DogTypeMst{}, err
	}
	return dogTypeMst, nil
}

// CreateDog: DBへdogのinsert
//
// args:
//   - model.Dog:	登録するdog
//
// return:
//   - model.dog:	登録されたdog
//   - error:	エラー
func (dr *dogRepository) CreateDog(c echo.Context, dog model.Dog) (model.Dog, error) {
	logger := log.GetLogger(c).Sugar()

	if err := dr.db.Create(&dog).Error; err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogのinsert処理で失敗しました。", errors.NewDogServerErrorEType())
		return model.Dog{}, err
	}
	return dog, nil
}

// UpdateDog: dogのupdate
//
// args:
//   - model.Dog:	更新するdog
//
// return:
//   - model.Dog:	更新したdog
//   - error:	エラー
func (dr *dogRepository) UpdateDog(c echo.Context, dog model.Dog) (model.Dog, error) {
	logger := log.GetLogger(c).Sugar()

	if err := dr.db.Save(&dog).Error; err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogのupdateで失敗しました。", errors.NewDogServerErrorEType())
		return model.Dog{}, err
	}
	return dog, nil
}

func (dr *dogRepository) DeleteDog(c echo.Context, dogID int64) error {
	logger := log.GetLogger(c).Sugar()

	result := dr.db.Where("dog_id=?", dogID).Delete(&model.Dog{})

	if result.Error != nil {
		logger.Error(result.Error)
		err := errors.NewWRError(result.Error, "dogのdelete処理で失敗しました。", errors.NewDogServerErrorEType())
		return err
	}
	if result.RowsAffected < 1 {
		logger.Error(result.Error)
		err := errors.NewWRError(result.Error, "dogのdelete処理で失敗しました。delete record is 0", errors.NewDogServerErrorEType())
		return err
	}
	return nil
}
