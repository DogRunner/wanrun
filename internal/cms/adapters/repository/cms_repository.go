package repository

import (
	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type ICmsRepository interface {
	CreateS3FileInfo(c echo.Context, s3FileInfo model.S3FileInfo) error
	GetS3FileInfoByFileID(c echo.Context, fileID string) ([]model.S3FileInfo, error)
	DeleteS3FileInfo(c echo.Context, s3FileInfo model.S3FileInfo) error
}

type cmsRepository struct {
	db *gorm.DB
}

func NewCmsRepository(db *gorm.DB) ICmsRepository {
	return &cmsRepository{db}
}

// CreateS3FileInfo: S3FileInfoの登録
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - model.S3FileInfo: S3FileInfoテーブルに登録する情報
//
// return:
//   - error: error情報
func (cr *cmsRepository) CreateS3FileInfo(c echo.Context, s3FileInfo model.S3FileInfo) error {
	logger := log.GetLogger(c).Sugar()

	// cmsテーブルにレコード作成
	if err := cr.db.Create(&s3FileInfo).Error; err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"DBへの登録が失敗しました。",
			wrErrors.NewCmsServerErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	return nil
}

// GetS3FileInfoByFileID: FileIDを元にS3FileInfo取得
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - string: fileID
//
// return:
//   - []model.S3FileInfo: S3ファイル情報
//   - error: error情報
func (cr *cmsRepository) GetS3FileInfoByFileID(c echo.Context, fileID string) ([]model.S3FileInfo, error) {
	logger := log.GetLogger(c).Sugar()

	s3Files := []model.S3FileInfo{}
	if err := cr.db.Model(&model.S3FileInfo{}).
		Where("file_id = ?", fileID).
		Find(&s3Files).Error; err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"DBからのデータ取得に失敗しました。",
			wrErrors.NewCmsServerErrorEType(),
		)
		logger.Errorf("DB search failure: %v", wrErr)

		return []model.S3FileInfo{}, wrErr
	}

	logger.Debugf("Query Results: %v", s3Files)

	return s3Files, nil
}

// DeleteS3FileInfo: S3FileInfoの削除
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - model.S3FileInfo: s3情報
//
// return:
//   - error: error情報
func (cr *cmsRepository) DeleteS3FileInfo(c echo.Context, s3Info model.S3FileInfo) error {
	logger := log.GetLogger(c).Sugar()

	if err := cr.db.Model(&model.S3FileInfo{}).
		Where("file_id = ? AND s3_object_key = ?", s3Info.FileID, s3Info.S3ObjectKey).
		Delete(&model.S3FileInfo{}).
		Error; err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"DBへの削除処理が失敗しました。",
			wrErrors.NewCmsServerErrorEType(),
		)
		logger.Errorf("DB delete failure: %v", wrErr)

		return wrErr
	}

	logger.Infof("Successfully deleted S3 file info with file_id: %s and s3_object_key: %s", s3Info.FileID.String, s3Info.S3ObjectKey.String)

	return nil
}
