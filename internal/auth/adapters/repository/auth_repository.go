package repository

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateDogOwner(c echo.Context, dogOC *model.DogOwnerCredential) (*model.DogOwnerCredential, error)
	// GetDogOwnerByEmail(c echo.Context, authDogOwner model.AuthDogOwner) (*model.AuthDogOwner, error)
	// CreateOAuthDogOwner(c echo.Context, dogOwnerCredential *model.DogOwnerCredential) (*model.DogOwnerCredential, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &authRepository{db}
}

func (ar *authRepository) CreateDogOwner(c echo.Context, dogOC *model.DogOwnerCredential) (*model.DogOwnerCredential, error) {
	logger := log.GetLogger(c).Sugar()

	// Emailの重複チェック
	if wrErr := ar.checkDuplicate(c, model.EmailField, dogOC.Email); wrErr != nil {
		return nil, wrErr
	}

	// PhoneNumberの重複チェック
	if wrErr := ar.checkDuplicate(c, model.PhoneNumberField, dogOC.PhoneNumber); wrErr != nil {
		return nil, wrErr
	}

	// トランザクションの開始
	err := ar.db.Transaction(func(tx *gorm.DB) error {
		// dog_ownersテーブルにレコード作成
		if err := tx.Create(&dogOC.AuthDogOwner.DogOwner).Error; err != nil {
			logger.Error("Failed to DogOwner: ", err)
			return err
		}

		// DogOwnerが作成された後、そのIDをauthDogOwnerに設定
		dogOC.AuthDogOwner.DogOwnerID = dogOC.AuthDogOwner.DogOwner.DogOwnerID

		// auth_dog_ownersテーブルにレコード作成
		if err := tx.Create(&dogOC.AuthDogOwner).Error; err != nil {
			logger.Error("Failed to AuthDogOwner: ", err)
			return err
		}
		// AuthDogOwnerが作成された後、そのIDをdogOwnerCredentialに設定
		dogOC.AuthDogOwnerID = dogOC.AuthDogOwner.AuthDogOwnerID

		// dog_owner_credentialsテーブルにレコード作成
		if err := tx.Create(&dogOC).Error; err != nil {
			logger.Error("Failed to DogOwnerCredential: ", err)
			return err
		}
		return nil
	})

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"DBへの登録が失敗しました。",
			wrErrors.NewDogownerServerErrorEType())

		logger.Errorf("Transaction failed error: %v", wrErr)

		return nil, wrErr
	}

	logger.Infof("Created DogOwner Detail: %v", dogOC.AuthDogOwner.DogOwner)
	logger.Infof("Created AuthDogOwner Detail: %v", dogOC.AuthDogOwner)
	logger.Infof("Created DogOwnerCredential Detail: %v", dogOC)

	return dogOC, nil
}

/*
OAuthユーザーの作成
*/
// func (ar *authRepository) CreateOAuthDogOwner(c echo.Context, dogOC *model.DogOwnerCredential) (*model.DogOwnerCredential, error) {
// 	logger := log.GetLogger(c).Sugar()

// 	// Emailの確認
// 	if wrErr := ar.checkDuplicate(c, model.PhoneNumberField, dogOC); wrErr != nil {
// 		return nil, wrErr
// 	}

// 	// PhoneNumberの確認

// 	// トランザクションの開始
// 	err := ar.db.Transaction(func(tx *gorm.DB) error {
// 		// dog_ownersテーブルにレコード作成
// 		if err := tx.Create(&dogOC.AuthDogOwner.DogOwner).Error; err != nil {
// 			logger.Error("Failed to DogOwner: ", err)
// 			return err
// 		}

// 		// DogOwnerが作成された後、そのIDをauthDogOwnerに設定
// 		dogOC.AuthDogOwner.DogOwnerID = dogOC.AuthDogOwner.DogOwner.DogOwnerID

// 		// auth_dog_ownersテーブルにレコード作成
// 		if err := tx.Create(&dogOC.AuthDogOwner).Error; err != nil {
// 			logger.Error("Failed to AuthDogOwner: ", err)
// 			return err
// 		}

// 		// AuthDogOwnerが作成された後、そのIDをdogOwnerCredentialに設定
// dogOC.AuthDogOwnerID = dogOC.AuthDogOwner.AuthDogOwnerID
// 		// dog_owner_credentialsテーブルにレコード作成
// 		if err := tx.Create(&dogOC).Error; err != nil {
// 			logger.Error("Failed to DogOwnerCredential: ", err)
// 			return err
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		wrErr := wrErrors.NewWRError(
// 			err,
// 			"transaction failed",
// 			wrErrors.NewDogownerClientErrorEType())

// 		logger.Errorf("Transaction failed error: %v", wrErr)

// 		return nil, wrErr
// 	}

// 	logger.Infof("Created DogOwner Detail: %v", dogOC.AuthDogOwner.DogOwner)
// 	logger.Infof("Created AuthDogOwner Detail: %v", dogOC.AuthDogOwner)
// 	logger.Infof("Created DogOwnerCredential Detail: %v", dogOC)

// 	// レスポンス用にDogOwnerのクレデンシャル取得
// 	var result model.DogOwnerCredential
// 	err = ar.db.Preload("AuthDogOwner").Preload("AuthDogOwner.DogOwner").First(&result, dogOC.CredentialID).Error

// 	if err != nil {
// 		wrErr := wrErrors.NewWRError(
// 			err,
// 			"failed to fetch created record",
// 			wrErrors.NewDogownerServerErrorEType())

// 		logger.Errorf("Failed to fetch created record: %v", wrErr)

// 		return nil, wrErr
// 	}

// 	logger.Infof("Created DogOwnerCredential Detail: %v", result)

// 	return &result, nil
// }

/*
DogOwnerのメールアドレスからDogOwner情報の取得
*/
// func (ar *authRepository) GetDogOwnerByEmail(c echo.Context, authDogOwner *model.AuthDogOwner) (*model.AuthDogOwner, error) {
// 	logger := log.GetLogger(c).Sugar()

// 	var result model.AuthDogOwner = model.AuthDogOwner{}

// 	if err := ar.db.Preload("DogOwner").Joins("LEFT JOIN dog_owners ON dog_owners.dog_owner_id = auth_dog_owners.dog_owner_id").Where("dog_owners.email = ?", authDogOwner.DogOwner.Email).First(&result).Error; err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	logger.Infof("Query Result: %v", result)

// 	return &result, nil
// }

/*
重複確認
*/
func (ar *authRepository) checkDuplicate(c echo.Context, field string, value sql.NullString) error {
	logger := log.GetLogger(c).Sugar()

	// 重複のvalidate
	var existingCount int64
	err := ar.db.Model(&model.DogOwnerCredential{}).Where(field+" = ?", value).Count(&existingCount).Error

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"DBからのデータ取得に失敗しました。",
			wrErrors.NewDogownerServerErrorEType(),
		)

		logger.Errorf("Failed to check existing value error: %v", wrErr)

		return wrErr
	}
	if existingCount > 0 {
		wrErr := wrErrors.NewWRError(
			nil,
			fmt.Sprintf("%sの%sが既に登録されています。", field, value.String),
			wrErrors.NewDogownerClientErrorEType(),
		)

		logger.Errorf("%s already exists error: %v", field, wrErr)

		return wrErr
	}
	return nil
}
