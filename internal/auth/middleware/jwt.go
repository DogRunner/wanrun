package middleware

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/configs"
	"github.com/wanrun-develop/wanrun/internal/auth/core/dto"
)

// CreateJwtConfig: context部分のconfig設定
//
// args:
//   - :
//   - :
//
// return:
//   - echo.MiddlewareFunc: ミドルウェア設定
func CreateJwtConfig(agp string) echo.MiddlewareFunc {
	return echojwt.WithConfig(
		echojwt.Config{
			SigningKey: []byte(configs.FetchCondigStr("jwt.os.secret.key")), // 署名用の秘密鍵
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return &dto.AccountClaims{} // カスタムクレームを設定
			},
			TokenLookup: "header:Authorization", // トークンの取得場所
			ContextKey:  "user_info",            // カスタムキーを設定
			Skipper: func(c echo.Context) bool {
				// authグループ配下のすべてのルートをスキップする
				return strings.HasPrefix(c.Path(), agp)
			},
		},
	)
}
