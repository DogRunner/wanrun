package model

import (
	"database/sql"
	"time"
)

type AuthDogrunmg struct {
	AuthDogrunmgID sql.NullInt64  `gorm:"primaryKey;column:auth_dogrun_manager_id;autoIncrement"`
	JwtID          sql.NullString `gorm:"size:45;column:jwt_id"`
	LoginAt        time.Time      `gorm:"column:login_at;not null;autoCreateTime"`

	Dogrunmg   Dogrunmg      `gorm:"foreignKey:DogrunmgID;references:DogrunmgID"`
	DogrunmgID sql.NullInt64 `gorm:"column:dogrun_manager_id;not null"`
}

type DogrunmgCredential struct {
	CredentialID sql.NullInt64  `gorm:"primaryKey;column:credential_id;autoIncrement"`
	Email        sql.NullString `gorm:"size:255;column:email"`
	Password     sql.NullString `gorm:"size:256;column:password"`
	IsAdmin      sql.NullBool   `gorm:"column:is_admin"`
	LoginAt      sql.NullTime   `gorm:"column:login_at;autoCreateTime"`

	AuthDogrunmg   AuthDogrunmg  `gorm:"foreignKey:AuthDogrunmgID;references:AuthDogrunmgID"`
	AuthDogrunmgID sql.NullInt64 `gorm:"column:auth_dogrun_manager_id;not null"`
}

func (AuthDogrunmg) TableName() string {
	return "auth_dogrun_managers"
}

func (DogrunmgCredential) TableName() string {
	return "dogrun_manager_credentials"
}
