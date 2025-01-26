package repository

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IDogrunmgScopeRepository interface {
	CreateDogrunmg(tx *gorm.DB, c echo.Context, adm *model.Dogrunmg) (sql.NullInt64, error)
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
//   - *model.Dogrunmg:dogrunmgの情報
//
// return:
//   - sql.NullInt64: dogrunmgのID
//   - error: error情報
func (dmsr *dogrunmgScopeRepository) CreateDogrunmg(
	tx *gorm.DB,
	c echo.Context,
	dm *model.Dogrunmg,
) (sql.NullInt64, error) {
	logger := log.GetLogger(c).Sugar()

	// Dogrunmgの作成
	if err := tx.Create(&dm).Error; err != nil {
		logger.Error("Failed to create Dogrunmg: ", err)
		return sql.NullInt64{}, wrErrors.NewWRError(
			err,
			"Dogrunmg作成に失敗しました。",
			wrErrors.NewDogrunmgServerErrorEType(),
		)
	}

	logger.Infof("Created Dogrunmg Detail: %v", dm)

	return dm.DogrunmgID, nil
}
