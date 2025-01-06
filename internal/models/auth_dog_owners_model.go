package model

import (
	"database/sql"
	"time"

	"github.com/wanrun-develop/wanrun/pkg/util"
)

const (
	OAUTH_GRANT_TYPE    string = "OAUTH"
	PASSWORD_GRANT_TYPE string = "PASSWORD"
)

type AuthDogowner struct {
	AuthDogownerID         sql.NullInt64   `gorm:"primaryKey;column:auth_dog_owner_id;autoIncrement"`
	AccessToken            sql.NullString  `gorm:"size:512;column:access_token"`
	RefreshToken           sql.NullString  `gorm:"size:512;column:refresh_token"`
	AccessTokenExpiration  util.CustomTime `gorm:"column:access_token_expiration"`
	RefreshTokenExpiration util.CustomTime `gorm:"column:refresh_token_expiration"`
	JwtID                  sql.NullString  `gorm:"size:45;column:jwt_id"`
	LoginAt                time.Time       `gorm:"column:login_at;not null;autoCreateTime"`

	Dogowner   Dogowner      `gorm:"foreignKey:DogownerID;references:DogownerID"`
	DogownerID sql.NullInt64 `gorm:"column:dog_owner_id;not null"`
}

type DogownerCredential struct {
	CredentialID sql.NullInt64 `gorm:"primaryKey;column:credential_id;autoIncrement"`
	// ProviderName   sql.NullString `gorm:"size:50;column:provider_name"`
	ProviderUserID sql.NullString `gorm:"size:256;column:provider_user_id"`
	Email          sql.NullString `gorm:"size:256;column:email"`
	PhoneNumber    sql.NullString `gorm:"size:15;column:phone_number"`
	Password       sql.NullString `gorm:"size:256;column:password"`
	GrantType      sql.NullString `gorm:"column:grant_type"`
	LoginAt        sql.NullTime   `gorm:"column:login_at;autoCreateTime"`

	AuthDogowner   AuthDogowner  `gorm:"foreignKey:AuthDogownerID;references:AuthDogownerID"`
	AuthDogownerID sql.NullInt64 `gorm:"column:auth_dog_owner_id;not null"`
}
