package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/configs"
	"github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"golang.org/x/exp/slices"
)

type IAuthJwt interface {
	NewJwtValidationMiddleware() echo.MiddlewareFunc
}

type authJwt struct {
	ar repository.IAuthRepository
}

func NewAuthJwt(ar repository.IAuthRepository) IAuthJwt {
	return &authJwt{ar}
}

const (
	CONTEXT_KEY   string = "user_info"
	TOKEN_LOOK_UP string = "header:Authorization:Bearer " // `Bearer `しか切り取れないのでスペースが多い場合は未対応
)

// スキップ対象のパスを定義
var skipPaths = []string{
	"/auth/token",
	"/auth/signUp",
	"/health",
}

// NewJwtValidationMiddleware: JWT検証用のミドルウェア設定を生成
//
// args:
//   - string: スキップするルートパスのプレフィックス（例: "/auth" で auth グループ配下のルートをスキップ）
//
// return:
//   - echo.MiddlewareFunc: JWT検証のためのミドルウェア設定
func (aj *authJwt) NewJwtValidationMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(
		echojwt.Config{
			SigningKey: []byte(configs.FetchConfigStr("jwt.os.secret.key")), // 署名用の秘密鍵
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return &handler.AccountClaims{} // カスタムクレームを設定
			},
			TokenLookup: TOKEN_LOOK_UP, // トークンの取得場所
			ContextKey:  CONTEXT_KEY,   // カスタムキーを設定
			Skipper: func(c echo.Context) bool { // スキップするパスを指定
				path := c.Path()
				return slices.Contains(skipPaths, path)
			},
			SuccessHandler: func(c echo.Context) {
				// contextからJWTのclaims取得と検証
				claims, wrErr := getJwtClaimsAndVerification(c)

				if wrErr != nil {
					neRes := errors.NewErrorRes(wrErr)
					_ = c.JSON(http.StatusUnauthorized, neRes)
					return
				}

				// リクエストのJWT内に含まれる`jwt_id`が、DBの`jwt_id`と一致するかを検証
				wrErr = aj.jwtIDValid(c, claims)

				if wrErr != nil {
					neRes := errors.NewErrorRes(wrErr)
					_ = c.JSON(http.StatusUnauthorized, neRes)
					return
				}

				// 全ての検証を終えたclaimsをcontextにセット
				c.Set(CONTEXT_KEY, claims)
			},
		},
	)
}

// getJwtClaimsAndVerification: contextからJWTのclaimsを取得と検証
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//
// return:
//   - *AccountClaims: contextから取得したJWTのクレーム情報
//   - error: error情報
func getJwtClaimsAndVerification(c echo.Context) (*handler.AccountClaims, error) {
	logger := log.GetLogger(c).Sugar()

	// JWTトークンをコンテキストから取得
	token, ok := c.Get(CONTEXT_KEY).(*jwt.Token)
	if !ok || token == nil {
		wrErr := errors.NewWRError(
			nil,
			"JWTトークンが見つかりません。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return nil, wrErr
	}

	// トークンが有効か確認
	if !token.Valid {
		wrErr := errors.NewWRError(
			nil,
			"無効なJWTトークンです。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return nil, wrErr
	}

	// クレーム情報を取得
	claims, ok := token.Claims.(*handler.AccountClaims)
	if !ok {
		wrErr := errors.NewWRError(
			nil,
			"クレーム情報の取得に失敗しました。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return nil, wrErr
	}

	// ExpiresAtの有効期限を確認
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		wrErr := errors.NewWRError(
			nil,
			"JWTトークンの有効期限が切れています。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return nil, wrErr
	}

	return claims, nil
}

// JwtValid: リクエストのJWT内に含まれる`jwt_id`が、DBの`jwt_id`と一致するかを検証
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *AccountClaims: contextから取得したJWTのクレーム情報
//
// return:
//   - error: error情報
func (aj *authJwt) jwtIDValid(c echo.Context, ac *handler.AccountClaims) error {
	logger := log.GetLogger(c).Sugar()

	// dogOwnerIDをstringからint64変換
	dogOwnerID, err := strconv.ParseInt(ac.ID, 10, 64)

	if err != nil {
		wrErr := errors.NewWRError(
			nil,
			"認証情報が違います。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	// 対象のdogOwnerからjwt_idの取得
	jwtID, wrErr := aj.ar.GetJwtID(c, dogOwnerID)

	if wrErr != nil {
		return wrErr
	}

	// フロントエンドからのリクエストのclaimsの`jti`とDBの`jwt_id`が一致するかどうか
	if jwtID != ac.JTI {
		wrErr := errors.NewWRError(
			nil,
			"jwt_idが一致しません。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	return nil
}
