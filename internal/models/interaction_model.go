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
