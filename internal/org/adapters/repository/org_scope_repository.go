package repository

import (
	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IOrgScopeRepository interface {
	CreateOrg(tx *gorm.DB, c echo.Context, drm *model.Dogrunmg) error
}

type orgScopeRepository struct {
}

func NewOrgScopeRepository() IOrgScopeRepository {
	return &orgScopeRepository{}
}

func (or *orgScopeRepository) CreateOrg(
	tx *gorm.DB,
	c echo.Context,
	drm *model.Dogrunmg) error {
	logger := log.GetLogger(c).Sugar()

	// organizationの作成
	if err := tx.Create(&drm.Organization).Error; err != nil {
		logger.Error("Failed to create Organization: ", err)
		return wrErrors.NewWRError(
			err,
			"Organization作成に失敗しました。",
			wrErrors.NewOrgServerErrorEType(),
		)
	}

	// Organizationが作成された後、そのIDをDogrun MGに設定
	drm.OrganizationID = drm.Organization.OrganizationID

	logger.Infof("Created Organization Detail: %v", drm.Organization)

	return nil
}
