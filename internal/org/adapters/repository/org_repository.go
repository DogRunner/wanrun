package repository

import (
	"gorm.io/gorm"
)

type IOrgRepository interface {
}

type orgRepository struct {
	db *gorm.DB
}

func NewOrgRepository(db *gorm.DB) IOrgRepository {
	return &orgRepository{
		db: db,
	}
}
