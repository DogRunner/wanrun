package repository

import (
	"errors"

	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IAuthJwtRepository interface {
	GetJwtID(c echo.Context, doi int64) (string, error)
}

type authJwtRepository struct {
	db *gorm.DB
}

// NewJwtRepository : JwtRepositoryのインスタンスを作成するコンストラクタ
func NewAuthJwtRepository(db *gorm.DB) IAuthJwtRepository {
	return &authJwtRepository{db}
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
func (ajr *authJwtRepository) GetJwtID(c echo.Context, doi int64) (string, error) {
	logger := log.GetLogger(c).Sugar()

	var result model.AuthDogOwner

	// 対象のdogOwnerのjwt_idの取得
	err := ajr.db.Model(&model.AuthDogOwner{}).
		Where("dog_owner_id= ?", doi).
		First(&result).
		Error

	if err != nil {
		// 空だった時
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wrErr := wrErrors.NewWRError(
				err,
				"認証情報がありません",
				wrErrors.NewDogownerClientErrorEType())

			logger.Errorf("Not found jwt id error: %v", wrErr)

			return "", wrErr
		}
		// その他のエラー処理
		wrErr := wrErrors.NewWRError(
			err,
			"DBからのデータ取得に失敗しました。",
			wrErrors.NewDogownerServerErrorEType())

		logger.Errorf("Failed to get JWT ID: %v", wrErr)

		return "", wrErr
	}
	logger.Debugf("Query Result: %v", result)

	return result.JwtID.String, nil
}
