package repository

import (
	"fmt"

	model "github.com/wanrun-develop/wanrun/internal/models"
	"gorm.io/gorm"
)

type IDogRepository interface {
	GetAllDogs() ([]model.Dog, error)
	GetDogByID(int64) (model.Dog, error)
	GetDogByDogOwnerID(int64) ([]model.Dog, error)
	CreateDog(model.Dog) (model.Dog, error)
	DeleteDog(int) error
}

type dogRepository struct {
	db *gorm.DB
}

func NewDogRepository(db *gorm.DB) IDogRepository {
	return &dogRepository{db}
}

func (dr *dogRepository) GetAllDogs() ([]model.Dog, error) {
	dogs := []model.Dog{}
	if err := dr.db.Find(&dogs).Error; err != nil {
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
func (dr *dogRepository) GetDogByID(dogID int64) (model.Dog, error) {
	dog := model.Dog{}
	if err := dr.db.Preload("DogType").Where("dog_id=?", dogID).First(&dog).Error; err != nil {
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
func (dr *dogRepository) GetDogByDogOwnerID(dogOwnerID int64) ([]model.Dog, error) {
	dogs := []model.Dog{}
	if err := dr.db.Preload("DogType").Where("dog_owner_id=?", dogOwnerID).Find(&dogs).Error; err != nil {
		return []model.Dog{}, err
	}
	return dogs, nil
}

// CreateDog: DBへdogのinsert
//
// args:
//   - model.Dog:	登録するdog
//
// return:
//   - model.dog:	登録されたdog
//   - error:	エラー
func (dr *dogRepository) CreateDog(dog model.Dog) (model.Dog, error) {
	if err := dr.db.Create(&dog).Error; err != nil {
		return model.Dog{}, err
	}
	return dog, nil
}

func (dr *dogRepository) DeleteDog(dogID int) error {
	result := dr.db.Where("dog_id=?", dogID).Delete(&model.Dog{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

// func (dr *dogRepository) UpdateDog(dogID uint) error {
// 	dog := model.Dog{}
// 	result := dr.db.Model(&dog).Clauses(clause.Returning{}).Where("dog_id=?", dogID).Update()
// }
