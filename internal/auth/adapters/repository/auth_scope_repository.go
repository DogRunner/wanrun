package repository

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IAuthScopeRepository interface {
	CreateAuthDogOwner(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error
	CreateDogOwnerCredential(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error
	CreateAuthDogrunmg(tx *gorm.DB, c echo.Context, adm *model.AuthDogrunmg) (sql.NullInt64, error)
	CreateDogrunmgCredential(tx *gorm.DB, c echo.Context, dmc *model.DogrunmgCredential) error
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
			wrErrors.NewAuthServerErrorEType(),
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
			wrErrors.NewAuthServerErrorEType(),
		)
	}

	logger.Infof("Created DogOwnerCredential Detail: %v", doc)

	return nil
}

// CreateAuthDogrunmg: AuthDogrunmgの登録処理
//
// args:
//   - *gorm.DB: トランザクションを張っているtx情報
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//   - *model.AuthDogrunmg: authのdogrunmgの情報
//
// return:
//   - sql.NullInt64:
//   - error: error情報
func (asr *authScopeRepository) CreateAuthDogrunmg(
	tx *gorm.DB,
	c echo.Context,
	adm *model.AuthDogrunmg,
) (sql.NullInt64, error) {
	logger := log.GetLogger(c).Sugar()

	// AuthDogrunmg作成
	if err := tx.Create(&adm).Error; err != nil {
		logger.Error("Failed to create AuthDogrunmg: ", err)
		return sql.NullInt64{}, wrErrors.NewWRError(
			err,
			"AuthDogrunmg作成に失敗しました。",
			wrErrors.NewAuthServerErrorEType(),
		)
	}

	logger.Infof("Created AuthDogrunmg Detail: %v", adm)

	return adm.AuthDogrunmgID, nil
}

// CreateDogrunmgCredential: DogrunmgのCredential登録処理
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//   - *model.DogrunmgCredential: dogrunmgのcredential情報
//
// return:
//   - error: error情報
func (asr *authScopeRepository) CreateDogrunmgCredential(
	tx *gorm.DB,
	c echo.Context,
	dmc *model.DogrunmgCredential,
) error {
	logger := log.GetLogger(c).Sugar()

	// DogrunmgのCredentials作成
	if err := tx.Create(&dmc).Error; err != nil {
		logger.Error("Failed to create DogrunmgCredential: ", err)
		return wrErrors.NewWRError(
			err,
			"DogrunmgCredential作成に失敗しました。",
			wrErrors.NewAuthServerErrorEType(),
		)
	}

	logger.Infof("Created DogrunmgCredential Detail: %v", dmc)

	return nil
}
