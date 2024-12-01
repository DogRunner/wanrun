package repository

import (
	model "github.com/wanrun-develop/wanrun/internal/models"
	"gorm.io/gorm"
)

type IDogOwnerRepository interface {
	GetDogOwnerById(int64) (model.DogOwner, error)
}

type dogOwnerRepository struct {
	db *gorm.DB
}

func NewDogRepository(db *gorm.DB) IDogOwnerRepository {
	return &dogOwnerRepository{db}
}

func (dr *dogOwnerRepository) GetDogOwnerById(dogOwnerId int64) (model.DogOwner, error) {
	dogOwner := model.DogOwner{}
	if err := dr.db.Where("dog_owner_id = ?", dogOwnerId).Find(&dogOwner).Error; err != nil {
		return model.DogOwner{}, err
	}
	return dogOwner, nil
}
