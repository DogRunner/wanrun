package facade

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
)

type IAuthFacade interface {
	OrgEmailValidate(c echo.Context, email string) error
}

type authFacade struct {
	ar repository.IAuthRepository
}

func NewAuthFacade(ar repository.IAuthRepository) IAuthFacade {
	return &authFacade{
		ar: ar,
	}
}

// OrgValidate: orgのEmailバリデーションフロー
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - string: 対象のEmail情報
//
// return:
//   - error: error情報
func (af *authFacade) OrgEmailValidate(c echo.Context, email string) error {
	// orgのEmail重複チェック
	if wrErr := af.ar.CheckOrgEmailExists(c, email); wrErr != nil {
		return wrErr
	}
	return nil
}
