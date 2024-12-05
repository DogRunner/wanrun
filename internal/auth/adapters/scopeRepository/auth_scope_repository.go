package scoperepository

import (
	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IAuthScopeRepository interface {
	CreateAuthDogOwner(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error
	CreateDogOwnerCredential(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error
}

type authScopeRepository struct {
}

func NewAuthScopeRepository() IAuthScopeRepository {
	return &authScopeRepository{}
}

// CreateAuthDogOwner: AuthDogOwnerの登録処理
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//   - *model.DogOwnerCredential: dogOwnerの情報
//
// return:
//   - error: error情報
func (asr *authScopeRepository) CreateAuthDogOwner(
	tx *gorm.DB,
	c echo.Context,
	doc *model.DogOwnerCredential,
) error {
	logger := log.GetLogger(c).Sugar()

	// auth_dog_ownersテーブルにAuthDogOwner作成
	if err := tx.Create(&doc.AuthDogOwner).Error; err != nil {
		logger.Error("Failed to create AuthDogOwner: ", err)
		return wrErrors.NewWRError(
			err,
			"AuthDogOwner作成に失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType(),
		)
	}

	// AuthDogOwnerが作成された後、そのIDをdogOwnerCredentialに設定
	doc.AuthDogOwnerID = doc.AuthDogOwner.AuthDogOwnerID

	logger.Infof("Created AuthDogOwner Detail: %v", doc.AuthDogOwner)

	return nil
}

// CreateDogOwnerCredential: DogOwnerのCredential登録処理
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//   - *model.DogOwnerCredential: dogOwnerの情報
//
// return:
//   - error: error情報
func (asr *authScopeRepository) CreateDogOwnerCredential(
	tx *gorm.DB,
	c echo.Context,
	doc *model.DogOwnerCredential,
) error {
	logger := log.GetLogger(c).Sugar()

	// dog_owner_credentialsテーブルにレコード作成
	if err := tx.Create(&doc).Error; err != nil {
		logger.Error("Failed to create DogOwnerCredential: ", err)
		return wrErrors.NewWRError(
			err,
			"DogOwnerCredential作成に失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType(),
		)
	}

	logger.Infof("Created DogOwnerCredential Detail: %v", doc)

	return nil
}
