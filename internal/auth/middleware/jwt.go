package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/configs"
	"github.com/wanrun-develop/wanrun/internal/auth/middleware/adapters/repository"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IAuthJwt interface {
	NewJwtValidationMiddleware() echo.MiddlewareFunc
}

type authJwt struct {
	ajr repository.IAuthJwtRepository
}

func NewAuthJwt(ajr repository.IAuthJwtRepository) IAuthJwt {
	return &authJwt{ajr}
}

const (
	RAW_CONTEXT_KEY      string = "user_info"
	TOKEN_LOOK_UP        string = "header:Authorization"
	VERIFIED_CONTEXT_KEY string = "claims"
)

// JWTのClaims
type AccountClaims struct {
	ID  string `json:"id"`
	JTI string `json:"jti"`
	jwt.RegisteredClaims
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
			SigningKey: []byte(configs.FetchCondigStr("jwt.os.secret.key")), // 署名用の秘密鍵
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return &AccountClaims{} // カスタムクレームを設定
			},
			TokenLookup: TOKEN_LOOK_UP,   // トークンの取得場所
			ContextKey:  RAW_CONTEXT_KEY, // カスタムキーを設定
			Skipper: func(c echo.Context) bool { // スキップするパスを指定
				path := c.Path()
				return path == "/auth/token" || path == "/auth/login"
			},
			SuccessHandler: func(c echo.Context) {
				// contextからJWTのclaims取得と検証
				claims, wrErr := getJwtClaimsAndVerification(c)

				if wrErr != nil {
					_ = c.JSON(http.StatusUnauthorized, map[string]error{"error": wrErr})
					return
				}

				// リクエストのJWT内に含まれる`jwt_id`が、DBの`jwt_id`と一致するかを検証
				isJwtIDValid, wrErr := aj.isJwtIDValid(c, claims)

				if !isJwtIDValid {
					_ = c.JSON(http.StatusUnauthorized, map[string]error{"error": wrErr})
					return
				}

				// 全ての検証を終えたclaimsをcontextにセット
				c.Set(VERIFIED_CONTEXT_KEY, claims)
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
func getJwtClaimsAndVerification(c echo.Context) (*AccountClaims, error) {
	logger := log.GetLogger(c).Sugar()

	var wrErr error
	// JWTトークンをコンテキストから取得
	token, ok := c.Get(RAW_CONTEXT_KEY).(*jwt.Token)
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
	claims, ok := token.Claims.(*AccountClaims)
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

	return claims, wrErr
}

// isJwtValid: リクエストのJWT内に含まれる`jwt_id`が、DBの`jwt_id`と一致するかを検証
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *AccountClaims: contextから取得したJWTのクレーム情報
//
// return:
//   - bool: jwt_idが一致したかどうかの結果
//   - error: error情報
func (aj *authJwt) isJwtIDValid(c echo.Context, ac *AccountClaims) (bool, error) {
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
		return false, err
	}

	// 対象のdogOwnerからjwt_idの取得
	jwtID, wrErr := aj.ajr.GetJwtID(c, dogOwnerID)

	if wrErr != nil {
		return false, wrErr
	}

	// フロントエンドからのリクエストのclaimsの`jti`とDBの`jwt_id`が一致するかどうか
	if jwtID != ac.JTI {
		wrErr := errors.NewWRError(
			nil,
			"jwt_idが一致しません。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return false, err
	}

	return true, nil
}
