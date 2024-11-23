package wrcontext

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	"github.com/wanrun-develop/wanrun/internal/auth/middleware"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

// GetClaims : contextから検証済みのclaims情報を取得する共通関数
//
// args:
//   - echo.MiddlewareFunc: JWT検証のためのミドルウェア設定
//
// return:
//   - *AccountClaims: 検証済みのclaims情報
func GetVerifiedClaims(c echo.Context) (*handler.AccountClaims, error) {
	logger := log.GetLogger(c).Sugar()

	claims, ok := c.Get(middleware.CONTEXT_KEY).(*handler.AccountClaims)
	if !ok || claims == nil {
		wrErr := errors.NewWRError(
			nil,
			"クレーム情報が見つかりません。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return nil, wrErr
	}

	return claims, nil
}

// GetLoginUserId: ログインユーザーIDの取得
//
//	コンテキストのjwt解析済みclaimからユーザーID取得
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - int64:	ユーザーID
func GetLoginUserId(c echo.Context) (int64, error) {
	logger := log.GetLogger(c).Sugar()

	claims, err := GetVerifiedClaims(c)
	if err != nil {
		return 0, err
	}
	userId, err := strconv.ParseInt(claims.ID, 10, 64)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(nil, "暫定", errors.NewAuthClientErrorEType())
		return 0, err
	}
	return userId, nil
}
