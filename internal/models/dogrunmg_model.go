package model

import (
	"database/sql"

	"github.com/wanrun-develop/wanrun/pkg/util"
)

type Dogrunmg struct {
	DogrunmgID sql.NullInt64   `gorm:"primaryKey;column:dogrun_manager_id;autoIncrement"`
	Name       sql.NullString  `gorm:"size:128;column:name;not null"`
	Image      sql.NullString  `json:"image" gorm:"type:text;column:image"`
	Sex        sql.NullString  `gorm:"size:1;column:sex"`
	CreateAt   util.CustomTime `gorm:"column:reg_at;not null;autoCreateTime"`
	UpdateAt   util.CustomTime `gorm:"column:upd_at;not null;autoUpdateTime"`

	// Orgとのリレーション
	Organization   Organization  `gorm:"foreignKey:OrganizationID;references:OrganizationID"`
	OrganizationID sql.NullInt64 `gorm:"column:organization_id;not null"` // 外部キー
}

func (Dogrunmg) TableName() string {
	return "dogrun_managers"
}

/*
Dogrunmgが空であるか
*/
func (dm *Dogrunmg) IsEmpty() bool {
	return !dm.IsNotEmpty()
}

/*
Dogrunmgが空でないか
*/
func (dm *Dogrunmg) IsNotEmpty() bool {
	return dm.DogrunmgID.Valid
}
