package repository

import (
	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IDogownerScopeRepository interface {
	CreateDogowner(tx *gorm.DB, c echo.Context, doc *model.DogownerCredential) error
}

type dogownerScopeRepository struct {
}

func NewDogownerScopeRepository() IDogownerScopeRepository {
	return &dogownerScopeRepository{}
}

// CreateDogowner: Dogownerの作成
//
// args:
//   - echo.Context: c Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *model.DogownerCredential: doc ドッグオーナーのクレデンシャル
//
// return:
//   - error: error情報
func (dosr *dogownerScopeRepository) CreateDogowner(tx *gorm.DB, c echo.Context, doc *model.DogownerCredential) error {
	logger := log.GetLogger(c).Sugar()

	// dog_ownersテーブルにDogownerの作成
	if err := tx.Create(&doc.AuthDogowner.Dogowner).Error; err != nil {
		logger.Error("Failed to create Dogowner: ", err)
		return wrErrors.NewWRError(
			err,
			"Dogowner作成に失敗しました。",
			wrErrors.NewDogownerServerErrorEType(),
		)
	}

	// Dogownerが作成された後、そのIDをauthDogownerに設定
	doc.AuthDogowner.DogownerID = doc.AuthDogowner.Dogowner.DogownerID

	logger.Infof("Created Dogowner Detail: %v", doc.AuthDogowner.Dogowner)

	return nil
}
