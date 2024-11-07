package controller

import (

	// "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/auth/core/dto"
	"github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IAuthController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	// GoogleOAuth(c echo.Context) error
}

type authController struct {
	ah handler.IAuthHandler
}

func NewAuthController(ah handler.IAuthHandler) IAuthController {
	return &authController{ah}
}

const (
	SIGNUP_MSG string = "飼い主の登録完了しました。"
	LOGIN_MSG  string = "ログインが完了しました。"
)

/*
GoogleのOAuth認証
*/
// func (ac *authController) GoogleOAuth(c echo.Context) error {
// 	logger := log.GetLogger(c).Sugar()

// 	// GrantTypeヘッダーを取得
// 	grantTypeHeader := c.Request().Header.Get(string(types.OAUTH_IDENTIFICATION_HEADER))

// 	// GrantTypeヘッダーのバリデーション
// 	if err := dto.ValidateGrantTypeHeader(grantTypeHeader, string(types.OAUTH_GRANT_TYPE_HEADER)); err != nil {
// 		err = wrErrors.NewWRError(err, "ヘッダー情報が異なります。", wrErrors.NewDogrunClientErrorEType())
// 		logger.Error(err)
// 		return err
// 	}

// 	// GrantTypeに型変換
// 	grantType := types.GrantType(grantTypeHeader)
// 	logger.Infof("grantTypeHeader: %v", grantType)

// 	// 認証コードの取得
// 	authorizationCode := c.QueryParam("code")

// 	// ユーザーが承認しなかったら、エラーのクエリパラメータにくるため
// 	oauthErrorCode := c.QueryParam("error")

// 	logger.Infof("authorizationCode: %v, oauthErrorCode: %v", authorizationCode, oauthErrorCode)

// 	// クエリパラメータのバリデーション
// 	if err := dto.ValidateOAuthResCode(authorizationCode, oauthErrorCode); err != nil {
// 		err = wrErrors.NewWRError(err, "承認をしてください。", wrErrors.NewDogrunClientErrorEType())
// 		logger.Error(err)
// 		return err
// 	}

// 	// OAuth認証
// 	resDogOwner, wrErr := ac.ah.GoogleOAuth(c, authorizationCode, grantType)

// 	if wrErr != nil {
// 		return wrErr
// 	}

// 	// jwt処理
// 	return jwtProcessing(c, resDogOwner)
// }

// SignUp: Password認証
//
// args:
//   - echo.Context: c Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//
// return:
//   - error: error情報
func (ac *authController) SignUp(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	reqADOD := dto.ReqAuthDogOwnerDto{}

	if err := c.Bind(&reqADOD); err != nil {
		wrErr := errors.NewWRError(err, "入力項目に不正があります。", errors.NewDogownerClientErrorEType())
		logger.Error(wrErr)
		return wrErr
	}

	// dogOwnerのSignUp
	resDogOwner, wrErr := ac.ah.SignUp(c, reqADOD)

	if wrErr != nil {
		return wrErr
	}

	// jwt処理
	return ac.ah.JwtProcessing(c, resDogOwner, SIGNUP_MSG)
}

// LogIn: login機能
//
// args:
//   - echo.Context: c Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//
// return:
//   - error: error情報
func (ac *authController) LogIn(c echo.Context) error {
	logger := log.GetLogger(c).Sugar()

	reqADOD := dto.ReqAuthDogOwnerDto{}

	if err := c.Bind(&reqADOD); err != nil {
		wrErr := errors.NewWRError(err, "入力項目に不正があります。", errors.NewDogownerClientErrorEType())
		logger.Error(wrErr)
		return wrErr
	}

	// LogIn処理
	resDogOwner, wrErr := ac.ah.LogIn(c, reqADOD)

	if wrErr != nil {
		return wrErr
	}

	// jwt処理
	return ac.ah.JwtProcessing(c, resDogOwner, LOGIN_MSG)
}

func (ac *authController) LogOut(c echo.Context) error { return nil }

// /*
// OAuthのクエリパラメータのバリデーション
// */
// func ValidateOAuthResCode(authorizationCode string, oauthErrorCode string) error {
// 	// "error" パラメータがある場合はエラーレスポンスを返す
// 	if oauthErrorCode != "" {
// 		wrErr := wrErrors.NewWRError(
// 			errOAuthFailed,
// 			"認証に失敗しました。",
// 			wrErrors.NewDogownerClientErrorEType(),
// 		)
// 		return wrErr
// 	}

// 	// "code" パラメータがある場合はそのまま処理
// 	if authorizationCode != "" {
// 		return nil
// 	}

// 	// どちらのパラメータもない場合は不正なリクエストとしてエラーを返す
// 	return errOAuthInvalidReq
// }
