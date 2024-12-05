package model

import (
	"database/sql"
	"time"
)

type S3FileInfo struct {
	S3FileInfoID sql.NullInt64  `gorm:"primaryKey;column:s3_file_info_id;autoIncrement"` // PK
	FileID       sql.NullString `gorm:"size:64;column:file_id;not null"`                 // 識別用UUID
	S3VersionID  sql.NullString `gorm:"size:256;column:s3_version_id"`                   // S3のバージョンID
	FileSize     sql.NullInt64  `gorm:"column:file_size;not null"`                       // ファイルサイズ
	S3ETag       sql.NullString `gorm:"size:64;column:s3_etag;not null"`                 // S3のETag
	S3ObjectKey  sql.NullString `gorm:"size:256;column:s3_object_key"`                   // S3のオブジェクトキー
	RegDate      time.Time      `gorm:"column:reg_date;autoCreateTime"`                  // 登録日時
	UpdDate      time.Time      `gorm:"column:upd_date;autoUpdateTime"`                  // 更新日時

	// DogOwnerとのリレーション
	DogOwner   DogOwner      `gorm:"foreignKey:DogOwnerID;references:DogOwnerID"`
	DogOwnerID sql.NullInt64 `gorm:"column:dog_owner_id;not null"` // dog_ownersのFK
}

func (S3FileInfo) TableName() string {
	return "s3_file_info" // 明示的にテーブル名を指定
}
