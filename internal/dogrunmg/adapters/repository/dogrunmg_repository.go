package repository

import (
	"gorm.io/gorm"
)

type IDogrunmgRepository interface {
}

type dogrunmgRepository struct {
	db *gorm.DB
}

func NewDogrunmgRepository(
	db *gorm.DB,
) IDogrunmgRepository {
	return &dogrunmgRepository{
		db: db,
	}
}
