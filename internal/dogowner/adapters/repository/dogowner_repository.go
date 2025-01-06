package repository

import (
	model "github.com/wanrun-develop/wanrun/internal/models"
	"gorm.io/gorm"
)

type IDogownerRepository interface {
	GetDogownerById(int64) (model.Dogowner, error)
}

type dogownerRepository struct {
	db *gorm.DB
}

func NewDogRepository(db *gorm.DB) IDogownerRepository {
	return &dogownerRepository{db}
}

func (dr *dogownerRepository) GetDogownerById(dogownerId int64) (model.Dogowner, error) {
	dogowner := model.Dogowner{}
	if err := dr.db.Where("dog_owner_id = ?", dogownerId).Find(&dogowner).Error; err != nil {
		return model.Dogowner{}, err
	}
	return dogowner, nil
}
