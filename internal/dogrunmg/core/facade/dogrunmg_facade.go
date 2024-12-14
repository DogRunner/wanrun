package facade

import (
	"github.com/labstack/echo/v4"
	dogrunmgRepository "github.com/wanrun-develop/wanrun/internal/dogrunmg/adapters/repository"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"gorm.io/gorm"
)

type IDogrunmgFacade interface {
	CreateOrg(tx *gorm.DB, c echo.Context, adm *model.AuthDogrunmg) error
}

type dogrunmgFacade struct {
	dmsr dogrunmgRepository.IDogrunmgScopeRepository
}

func NewDogrunmgFacade(
	dmsr dogrunmgRepository.IDogrunmgScopeRepository,
) IDogrunmgFacade {
	return &dogrunmgFacade{
		dmsr: dmsr,
	}
}

// CreateOrg: dogrunmg登録フロー
//
// args:
//   - *gorm.DB: トランザクションを張っているtx情報
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *model.AuthDogrunmg: authdogrunmgの情報
//
// return:
//   - error: error情報
func (dmf *dogrunmgFacade) CreateOrg(
	tx *gorm.DB,
	c echo.Context,
	adm *model.AuthDogrunmg,
) error {

	if wrErr := dmf.dmsr.CreateDogrunmg(tx, c, adm); wrErr != nil {
		return wrErr
	}
	return nil
}
