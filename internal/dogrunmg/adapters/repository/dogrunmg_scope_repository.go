package repository

import (
	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IDogrunmgScopeRepository interface {
	CreateDogrunmg(tx *gorm.DB, c echo.Context, adm *model.AuthDogrunmg) error
}

type dogrunmgScopeRepository struct {
}

func NewDogrunmgScopeRepository() IDogrunmgScopeRepository {
	return &dogrunmgScopeRepository{}
}

// CreateDogrunmg: Dogrunmgの作成
//
// args:
//   - *gorm.DB: トランザクションを張っているtx情報
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *model.AuthDogrunmg: authのdogrunmgの情報
//
// return:
//   - error: error情報
func (dmsr *dogrunmgScopeRepository) CreateDogrunmg(
	tx *gorm.DB,
	c echo.Context,
	adm *model.AuthDogrunmg,
) error {
	logger := log.GetLogger(c).Sugar()

	// Dogrunmgの作成
	if err := tx.Create(&adm.Dogrunmg).Error; err != nil {
		logger.Error("Failed to create Dogrunmg: ", err)
		return wrErrors.NewWRError(
			err,
			"Dogrunmg作成に失敗しました。",
			wrErrors.NewDogrunmgServerErrorEType(),
		)
	}

	// Dogrunmgが作成された後、そのIDをauthDogrunmgに設定
	adm.DogrunmgID = adm.Dogrunmg.DogrunmgID

	logger.Infof("Created Dogrunmg Detail: %v", adm.Dogrunmg)

	return nil
}
