package model

import (
	"database/sql"
)

type Dog struct {
	DogID      sql.NullInt64  `gorm:"primaryKey;column:dog_id;autoIncrement"`
	DogOwnerID sql.NullInt64  `gorm:"column:dog_owner_id;not null;foreignKey:DogOwnerID"`
	Name       sql.NullString `gorm:"size:128;column:name;not null"`
	DogTypeID  sql.NullInt64  `gorm:"column:dog_type_id"`
	Weight     sql.NullInt64  `gorm:"column:weight"`
	Sex        sql.NullString `gorm:"size:1;column:sex"`
	Image      sql.NullString `gorm:"column:image"`
	CreateAt   sql.NullTime   `gorm:"column:reg_at;not null;autoCreateTime"`
	UpdateAt   sql.NullTime   `gorm:"column:upd_at;not null;autoUpdateTime"`

	//リレーション
	DogOwner DogOwner   `gorm:"foreignKey:DogOwnerID;references:DogOwnerID"`
	DogType  DogTypeMst `gorm:"foreignKey:DogTypeID;references:DogTypeID"`
}

// dogが空かの判定
func (d *Dog) IsEmpty() bool {
	return !d.DogID.Valid
}

type DogTypeMst struct {
	DogTypeID int    `gorm:"primaryKey;column:dog_type_id"`
	Name      string `gorm:"column:name;not null"`
}

// GORMにテーブル名を指定
func (DogTypeMst) TableName() string {
	return "dog_type_mst"
}
