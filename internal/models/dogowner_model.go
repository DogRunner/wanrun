package model

import (
	"database/sql"

	"github.com/wanrun-develop/wanrun/pkg/util"
)

type Dogowner struct {
	DogownerID sql.NullInt64   `json:"dogOwnerId" gorm:"primaryKey;column:dog_owner_id;autoIncrement"`
	Name       sql.NullString  `json:"name" gorm:"size:128;column:name;not null"`
	Image      sql.NullString  `json:"image" gorm:"type:text;column:image"`
	Sex        sql.NullString  `json:"sex" gorm:"size:1;column:sex"`
	CreateAt   util.CustomTime `json:"createAt" gorm:"column:reg_at;not null;autoCreateTime"`
	UpdateAt   util.CustomTime `json:"updateAt" gorm:"column:upd_at;not null;autoCreateTime"`
}

// dogownerが空かの判定
func (do *Dogowner) IsEmpty() bool {
	return !do.DogownerID.Valid
}
