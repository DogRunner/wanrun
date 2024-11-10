package internal

import (
	"os"

	"github.com/labstack/echo/v4"
	authRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/auth/middleware"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

func Test(c echo.Context, dbConn *gorm.DB) error {
	logger := log.GetLogger(c).Sugar()
	authRepository := authRepository.NewAuthRepository(dbConn)
	authJwt := middleware.NewAuthJwt(authRepository)

	claims, wrErr := middleware.GetJwtClaims(c)

	if wrErr != nil {
		return wrErr
	}

	authJwt.IsJwtIDValid(c, claims)

	logger.Info("Test*()の実行. ")
	if err := testError(); err != nil {
		err = errors.NewWRError(err, "エラー再生成しました。", errors.NewAuthClientErrorEType())
		logger.Error(err)
		return err
	}
	return nil
}

func testError() error {
	file := "xxx/xxx"
	_, err := os.Open(file)
	if err != nil {
		err := errors.NewWRError(err, "エラー発生: entityFuncのファイル読み込み", errors.NewAuthClientErrorEType())
		return err
	}
	return nil
}
