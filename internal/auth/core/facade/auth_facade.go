package facade

import (
	"github.com/labstack/echo/v4"
	authScopeRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/scoperepository"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"gorm.io/gorm"
)

type IAuthFacade interface {
	CreateOrg(tx *gorm.DB, c echo.Context, dmc *model.DogrunmgCredential) error
}

type authFacade struct {
	asr authScopeRepository.IAuthScopeRepository
}

func NewAuthFacade(
	asr authScopeRepository.IAuthScopeRepository,
) IAuthFacade {
	return &authFacade{
		asr: asr,
	}
}

// CreateOrg: authDogrunmgとdogrunmgのCredential登録フロー
//
// args:
//   - *gorm.DB: トランザクション情報
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//   - *model.DogrunmgCredential: dogrunmgの情報
//
// return:
//   - error: error情報
func (af *authFacade) CreateOrg(
	tx *gorm.DB,
	c echo.Context,
	dmc *model.DogrunmgCredential,
) error {

	// AuthDogrunmgの作成
	if wrErr := af.asr.CreateAuthDogrunmg(tx, c, dmc); wrErr != nil {
		return nil
	}

	// DogrunmgのCredentialsの作成
	if wrErr := af.asr.CreateDogrunmgCredential(tx, c, dmc); wrErr != nil {
		return nil
	}

	return nil
}
