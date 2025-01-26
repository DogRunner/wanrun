package facade

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
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

// OrgValidate: OrgのEmailバリデーションフロー
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - string: 対象のEmail情報
//
// return:
//   - error: error情報
func (af *authFacade) OrgEmailValidate(c echo.Context, email string) error {
	logger := log.GetLogger(c).Sugar()

	// OrgのEmail数取得
	existingCount, wrErr := af.ar.CountOrgEmail(c, email)

	if wrErr != nil {
		return wrErr
	}

	// Emailの重複確認
	if existingCount > 0 {
		wrErr := wrErrors.NewWRError(
			nil,
			fmt.Sprintf("%sのEmailが既に登録されています。", email),
			wrErrors.NewDogOwnerClientErrorEType(),
		)

		logger.Errorf("%s already exists error: %v", email, wrErr)

		return wrErr
	}

	return nil
}
