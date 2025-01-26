package repository

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IOrgScopeRepository interface {
	CreateOrg(tx *gorm.DB, c echo.Context, o *model.Organization) (sql.NullInt64, error)
}

type orgScopeRepository struct {
}

func NewOrgScopeRepository() IOrgScopeRepository {
	return &orgScopeRepository{}
}

// CreateOrg: organizationの作成
//
// args:
//   - *gorm.DB: トランザクションを張っているtx情報
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *model.Organization: organization情報
//
// return:
//   - sql.NullInt64: organizationのID
//   - error: error情報
func (or *orgScopeRepository) CreateOrg(
	tx *gorm.DB,
	c echo.Context,
	o *model.Organization,
) (sql.NullInt64, error) {
	logger := log.GetLogger(c).Sugar()

	// organizationの作成
	if err := tx.Create(&o).Error; err != nil {
		logger.Error("Failed to create Organization: ", err)
		return sql.NullInt64{}, wrErrors.NewWRError(
			err,
			"Organization作成に失敗しました。",
			wrErrors.NewOrgServerErrorEType(),
		)
	}

	logger.Infof("Created Organization Detail: %v", o)

	return o.OrganizationID, nil
}
