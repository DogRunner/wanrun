package model

import "database/sql"

type DogrunBookmark struct {
	DogrunBookmarkID sql.NullInt64 `gorm:"column:dogrun_bookmark_id;primaryKey"`
	DogOwnerID       sql.NullInt64 `gorm:"column:dog_owner_id;not null"`
	DogrunID         sql.NullInt64 `gorm:"column:dogrun_id;not null"`
	SavedAt          sql.NullTime  `gorm:"column:saved_at;autoCreateTime"`
}

/*
DogrunBookmarkが空であるか
*/
func (b *DogrunBookmark) IsEmpty() bool {
	return !b.IsNotEmpty()
}

/*
DogrunBookmarkが空でないか
*/
func (b *DogrunBookmark) IsNotEmpty() bool {
	return b.DogrunBookmarkID.Valid
}

type DogrunCheckin struct {
	DogrunCheckinID sql.NullInt64 `gorm:"column:dogrun_checkin_id;primaryKey"`
	DogrunID        sql.NullInt64 `gorm:"column:dogrun_id;not null"`
	DogID           sql.NullInt64 `gorm:"column:dog_id;not null"`
	CheckinAt       sql.NullTime  `gorm:"column:checkin_at;autoCreateTime"`
	ReCheckinAt     sql.NullTime  `gorm:"column:re_checkin_at;autoUpdateTime"`
}

func (DogrunCheckin) TableName() string {
	return "dogrun_checkin"
}

/*
DogrunCheckinが空であるか
*/
func (ci *DogrunCheckin) IsEmpty() bool {
	return !ci.IsNotEmpty()
}

/*
DogrunCheckinが空でないか
*/
func (ci *DogrunCheckin) IsNotEmpty() bool {
	return ci.DogrunCheckinID.Valid
}

type DogrunCheckout struct {
	DogrunCheckoutID sql.NullInt64 `gorm:"column:dogrun_checkout_id;primaryKey"`
	DogrunID         sql.NullInt64 `gorm:"column:dogrun_id;not null"`
	DogID            sql.NullInt64 `gorm:"column:dog_id;not null"`
	CheckoutAt       sql.NullTime  `gorm:"column:checkout_at;autoCreateTime"`
	ReCheckoutAt     sql.NullTime  `gorm:"column:re_checkout_at;autoUpdateTime"`
}

func (DogrunCheckout) TableName() string {
	return "dogrun_checkout"
}

/*
DogrunCheckinが空であるか
*/
func (co *DogrunCheckout) IsEmpty() bool {
	return !co.IsNotEmpty()
}

/*
DogrunCheckinが空でないか
*/
func (co *DogrunCheckout) IsNotEmpty() bool {
	return co.DogrunCheckoutID.Valid
}
