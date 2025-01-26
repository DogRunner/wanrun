package model

import (
	"database/sql"

	"github.com/wanrun-develop/wanrun/pkg/util"
)

type Organization struct {
	OrganizationID sql.NullInt64   `gorm:"primaryKey;column:organization_id;autoIncrement"`
	Name           sql.NullString  `gorm:"size:128;column:organization_name;not null"`
	ContactEmail   sql.NullString  `gorm:"size:256;column:contact_email"`
	PhoneNumber    sql.NullString  `gorm:"size:15;column:phone_number"`
	Address        sql.NullString  `gorm:"size:256;column:address"`
	Description    sql.NullString  `gorm:"size:512;column:description"`
	CreateAt       util.CustomTime `gorm:"column:reg_at;not null;autoCreateTime"`
	UpdateAt       util.CustomTime `gorm:"column:upd_at;not null;autoCreateTime"`
}

// dogownerが空かの判定
func (o *Organization) IsEmpty() bool {
	return !o.OrganizationID.Valid
}
