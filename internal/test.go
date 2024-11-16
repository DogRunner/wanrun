package internal

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/wrcontext"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

func Test(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	// claims情報の取得
	claims, wrErr := wrcontext.GetVerifiedClaims(c)

	if wrErr != nil {
		return wrErr
	}

	userID := claims.ID
	jti := claims.JTI
	exp := claims.ExpiresAt

	logger.Infof("userID: %v, jti: %v, exp: %v\n", userID, jti, exp)

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
