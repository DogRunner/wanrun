package repository

import (
	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IDogOwnerScopeRepository interface {
	CreateDogOwner(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error
}

type dogOwnerScopeRepository struct {
}

func NewDogOwnerScopeRepository() IDogOwnerScopeRepository {
	return &dogOwnerScopeRepository{}
}

// CreateDogOwner: DogOwnerの作成
//
// args:
//   - echo.Context: c Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *model.DogOwnerCredential: doc ドッグオーナーのクレデンシャル
//
// return:
//   - error: error情報
func (dosr *dogOwnerScopeRepository) CreateDogOwner(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error {
	logger := log.GetLogger(c).Sugar()

	// dog_ownersテーブルにDogOwnerの作成
	if err := tx.Create(&doc.AuthDogOwner.DogOwner).Error; err != nil {
		logger.Error("Failed to create DogOwner: ", err)
		return wrErrors.NewWRError(
			err,
			"DogOwner作成に失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType(),
		)
	}

	// DogOwnerが作成された後、そのIDをauthDogOwnerに設定
	doc.AuthDogOwner.DogOwnerID = doc.AuthDogOwner.DogOwner.DogOwnerID

	logger.Infof("Created DogOwner Detail: %v", doc.AuthDogOwner.DogOwner)

	return nil
}
