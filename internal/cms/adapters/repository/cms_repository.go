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

	// トランザクションの開始
	err := cr.db.Transaction(func(tx *gorm.DB) error {
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
	})

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"トランザクションが失敗しました。",
			wrErrors.NewCmsServerErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	return nil
}
