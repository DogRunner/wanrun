package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/auth/core/dto"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateDogOwner(c echo.Context, doc *model.DogOwnerCredential) (*model.DogOwnerCredential, error)
	GetDogOwnerByCredential(c echo.Context, ador dto.AuthDogOwnerReq) (*model.DogOwnerCredential, error)
	// CreateOAuthDogOwner(c echo.Context, dogOwnerCredential *model.DogOwnerCredential) (*model.DogOwnerCredential, error)
	UpdateJwtID(c echo.Context, doc *model.DogOwnerCredential, jwt_id string) error
	GetJwtID(c echo.Context, doi int64) (string, error)
	DeleteJwtID(c echo.Context, doID int64) error
	CheckDuplicate(c echo.Context, field string, value sql.NullString) error
	CreateAuthDogOwnerAndCredential(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &authRepository{db}
}

// CreateDogOwner: DogOwnerの作成
//
// args:
//   - echo.Context: c Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *model.DogOwnerCredential: doc ドッグオーナーのクレデンシャル
//
// return:
//   - *model.DogOwnerCredential: ドッグオーナーのクレデンシャル
//   - error: error情報
func (ar *authRepository) CreateDogOwner(c echo.Context, doc *model.DogOwnerCredential) (*model.DogOwnerCredential, error) {
	logger := log.GetLogger(c).Sugar()

	// Emailの重複チェック
	if wrErr := ar.CheckDuplicate(c, model.EmailField, doc.Email); wrErr != nil {
		return nil, wrErr
	}

	// PhoneNumberの重複チェック
	if wrErr := ar.CheckDuplicate(c, model.PhoneNumberField, doc.PhoneNumber); wrErr != nil {
		return nil, wrErr
	}

	// トランザクションの開始
	err := ar.db.Transaction(func(tx *gorm.DB) error {
		// dog_ownersテーブルにレコード作成
		if err := tx.Create(&doc.AuthDogOwner.DogOwner).Error; err != nil {
			logger.Error("Failed to DogOwner: ", err)
			return err
		}

		// DogOwnerが作成された後、そのIDをauthDogOwnerに設定
		doc.AuthDogOwner.DogOwnerID = doc.AuthDogOwner.DogOwner.DogOwnerID

		// auth_dog_ownersテーブルにレコード作成
		if err := tx.Create(&doc.AuthDogOwner).Error; err != nil {
			logger.Error("Failed to AuthDogOwner: ", err)
			return err
		}
		// AuthDogOwnerが作成された後、そのIDをdogOwnerCredentialに設定
		doc.AuthDogOwnerID = doc.AuthDogOwner.AuthDogOwnerID

		// dog_owner_credentialsテーブルにレコード作成
		if err := tx.Create(&doc).Error; err != nil {
			logger.Error("Failed to DogOwnerCredential: ", err)
			return err
		}
		return nil
	})

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"DBへの登録が失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType())

		logger.Errorf("Transaction failed error: %v", wrErr)

		return nil, wrErr
	}

	logger.Infof("Created DogOwner Detail: %v", doc.AuthDogOwner.DogOwner)
	logger.Infof("Created AuthDogOwner Detail: %v", doc.AuthDogOwner)
	logger.Infof("Created DogOwnerCredential Detail: %v", doc)

	return doc, nil
}

/*
OAuthユーザーの作成
*/
// func (ar *authRepository) CreateOAuthDogOwner(c echo.Context, doc *model.DogOwnerCredential) (*model.DogOwnerCredential, error) {
// 	logger := log.GetLogger(c).Sugar()

// 	// Emailの確認
// 	if wrErr := ar.checkDuplicate(c, model.PhoneNumberField, doc); wrErr != nil {
// 		return nil, wrErr
// 	}

// 	// PhoneNumberの確認

// 	// トランザクションの開始
// 	err := ar.db.Transaction(func(tx *gorm.DB) error {
// 		// dog_ownersテーブルにレコード作成
// 		if err := tx.Create(&doc.AuthDogOwner.DogOwner).Error; err != nil {
// 			logger.Error("Failed to DogOwner: ", err)
// 			return err
// 		}

// 		// DogOwnerが作成された後、そのIDをauthDogOwnerに設定
// 		doc.AuthDogOwner.DogOwnerID = doc.AuthDogOwner.DogOwner.DogOwnerID

// 		// auth_dog_ownersテーブルにレコード作成
// 		if err := tx.Create(&doc.AuthDogOwner).Error; err != nil {
// 			logger.Error("Failed to AuthDogOwner: ", err)
// 			return err
// 		}

// 		// AuthDogOwnerが作成された後、そのIDをdogOwnerCredentialに設定
// doc.AuthDogOwnerID = doc.AuthDogOwner.AuthDogOwnerID
// 		// dog_owner_credentialsテーブルにレコード作成
// 		if err := tx.Create(&doc).Error; err != nil {
// 			logger.Error("Failed to DogOwnerCredential: ", err)
// 			return err
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		wrErr := wrErrors.NewWRError(
// 			err,
// 			"transaction failed",
// 			wrErrors.NewDogOwnerClientErrorEType())

// 		logger.Errorf("Transaction failed error: %v", wrErr)

// 		return nil, wrErr
// 	}

// 	logger.Infof("Created DogOwner Detail: %v", doc.AuthDogOwner.DogOwner)
// 	logger.Infof("Created AuthDogOwner Detail: %v", doc.AuthDogOwner)
// 	logger.Infof("Created DogOwnerCredential Detail: %v", doc)

// 	// レスポンス用にDogOwnerのクレデンシャル取得
// 	var result model.DogOwnerCredential
// 	err = ar.db.Preload("AuthDogOwner").Preload("AuthDogOwner.DogOwner").First(&result, doc.CredentialID).Error

// 	if err != nil {
// 		wrErr := wrErrors.NewWRError(
// 			err,
// 			"failed to fetch created record",
// 			wrErrors.NewDogOwnerServerErrorEType())

// 		logger.Errorf("Failed to fetch created record: %v", wrErr)

// 		return nil, wrErr
// 	}

// 	logger.Infof("Created DogOwnerCredential Detail: %v", result)

// 	return &result, nil
// }

// GetDogOwnerByCredential: ドッグオーナーのクレデンシャル取得
//
// args:
//   - echo.Context: c Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - dto.AuthDogOwnerReq: authDogOwnerのリクエスト情報
//
// return:
//   - *model.DogOwnerCredential: ドッグオーナーのクレデンシャル
//   - error: error情報
func (ar *authRepository) GetDogOwnerByCredential(c echo.Context, ador dto.AuthDogOwnerReq) (*model.DogOwnerCredential, error) {
	logger := log.GetLogger(c).Sugar()

	var result model.DogOwnerCredential

	// EmailまたはPhoneNumberとgrantTypeがPASSWORDに基づくレコードを検索
	err := ar.db.Model(&model.DogOwnerCredential{}).
		Where("(email = ? OR phone_number = ?) AND grant_type = ?", ador.Email, ador.PhoneNumber, model.PASSWORD_GRANT_TYPE).
		Preload("AuthDogOwner").
		Preload("AuthDogOwner.DogOwner").
		First(&result).
		Error

	if err != nil {
		// 空だった時
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wrErr := wrErrors.NewWRError(
				err,
				"認証情報がありません",
				wrErrors.NewDogOwnerClientErrorEType())

			logger.Errorf("Not found credential error: %v", wrErr)

			return nil, wrErr
		}
		// その他のエラー処理
		wrErr := wrErrors.NewWRError(
			err,
			"DBからのデータ取得に失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType())

		logger.Errorf("DB search failure: %v", wrErr)

		return nil, wrErr
	}
	logger.Debugf("Query Result: %v", result)

	return &result, nil
}

// UpdateJwtID: 対象のdogOwnerのjwt_idの更新
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *model.DogOwnerCredential: dogOwnerの情報
//   - string: 更新用のjwt_id
//
// return:
//   - error: error情報
func (ar *authRepository) UpdateJwtID(c echo.Context, doc *model.DogOwnerCredential, ji string) error {
	logger := log.GetLogger(c).Sugar()

	// 対象のdogOwnerのjwt_idの更新
	err := ar.db.Model(&model.AuthDogOwner{}).
		Where("dog_owner_id= ?", doc.AuthDogOwner.DogOwner.DogOwnerID.Int64).
		Update("jwt_id", ji).Error

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"DBへの更新が失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType())

		logger.Errorf("Failed to update JWT ID: %v", wrErr)

		return wrErr
	}

	return err
}

// DeleteJwtID: 対象のdogOwnerのjwt_idの削除
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - dto.DogOwnerDTO: dogOwnerの情報
//
// return:
//   - error: error情報
func (ar *authRepository) DeleteJwtID(c echo.Context, doID int64) error {
	logger := log.GetLogger(c).Sugar()

	// 対象のdogOwnerのjwt_idの更新
	err := ar.db.Model(&model.AuthDogOwner{}).
		Where("dog_owner_id= ?", doID).
		Update("jwt_id", nil).Error
	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"DBへの同期が失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType())

		logger.Errorf("Failed to delete JWT ID: %v", wrErr)

		return wrErr
	}

	return err
}

// GetJwtID:
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - int64: 取得したいdogOwnerID
//
// return:
//   - string: 対象のdogOwnerのjwt_id
//   - error: error情報
func (ar *authRepository) GetJwtID(c echo.Context, doi int64) (string, error) {
	logger := log.GetLogger(c).Sugar()

	var result model.AuthDogOwner

	// 対象のdogOwnerのjwt_idの取得
	err := ar.db.Model(&model.AuthDogOwner{}).
		Where("dog_owner_id= ?", doi).
		First(&result).
		Error

	if err != nil {
		// 空だった時
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wrErr := wrErrors.NewWRError(
				err,
				"認証情報がありません",
				wrErrors.NewDogOwnerClientErrorEType())

			logger.Errorf("Not found jwt id error: %v", wrErr)

			return "", wrErr
		}
		// その他のエラー処理
		wrErr := wrErrors.NewWRError(
			err,
			"DBからのデータ取得に失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType())

		logger.Errorf("Failed to get JWT ID: %v", wrErr)

		return "", wrErr
	}
	logger.Debugf("Query Result: %v", result)

	return result.JwtID.String, nil
}

// checkDuplicate:  Password認証のバリデーション
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - string: 対象のdbのフィールド名
//   - sql.NullString: 対象のgrant_typeの値
//
// return:
//   - error: error情報
func (ar *authRepository) CheckDuplicate(c echo.Context, field string, value sql.NullString) error {
	logger := log.GetLogger(c).Sugar()

	// 重複のvalidate
	var existingCount int64

	// grant_typeがPASSWORDで重複していないかの確認
	err := ar.db.Model(&model.DogOwnerCredential{}).
		Where(field+" = ? AND grant_type = ?", value, model.PASSWORD_GRANT_TYPE).
		Count(&existingCount).
		Error

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"DBからのデータ取得に失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType(),
		)

		logger.Errorf("Failed to check existing value error: %v", wrErr)

		return wrErr
	}

	if existingCount > 0 {
		wrErr := wrErrors.NewWRError(
			nil,
			fmt.Sprintf("%sの%sが既に登録されています。", field, value.String),
			wrErrors.NewDogOwnerClientErrorEType(),
		)

		logger.Errorf("%s already exists error: %v", field, wrErr)

		return wrErr
	}
	return nil
}

func (ar *authRepository) CreateAuthDogOwnerAndCredential(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error {
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

	// dog_owner_credentialsテーブルにレコード作成
	if err := tx.Create(&doc).Error; err != nil {
		logger.Error("Failed to create DogOwnerCredential: ", err)
		return wrErrors.NewWRError(
			err,
			"DogOwnerCredential作成に失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType(),
		)
	}

	logger.Infof("Created AuthDogOwner Detail: %v", doc.AuthDogOwner)
	logger.Infof("Created DogOwnerCredential Detail: %v", doc)

	return nil
}
