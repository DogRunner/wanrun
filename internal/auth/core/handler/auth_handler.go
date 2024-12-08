package handler

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/configs"
	"github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	authDTO "github.com/wanrun-develop/wanrun/internal/auth/core/dto"
	doDTO "github.com/wanrun-develop/wanrun/internal/dogowner/core/dto"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

type IAuthHandler interface {
	LogIn(c echo.Context, ador authDTO.AuthDogOwnerReq) (string, error)
	Revoke(c echo.Context, claims *AccountClaims) error
	// GoogleOAuth(c echo.Context, authorizationCode string, grantType types.GrantType) (dto.ResDogOwnerDto, error)
}

type authHandler struct {
	ar repository.IAuthRepository
	// ag google.IOAuthGoogle
}

//	func NewAuthHandler(ar repository.IAuthRepository, g google.IOAuthGoogle) IAuthHandler {
//		return &authHandler{ar, g}
//	}
func NewAuthHandler(ar repository.IAuthRepository) IAuthHandler {
	return &authHandler{ar}
}

// JWTのClaims
type AccountClaims struct {
	ID  string `json:"id"`
	JTI string `json:"jti"`
	jwt.RegisteredClaims
}

// GetDogOwnerIDAsInt64: 共通処理で、int64のDogOwnerのID取得。
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//
// return:
//   - int64: dogOwnerのID情報(int64)
//   - error: error情報
func (claims *AccountClaims) GetDogOwnerIDAsInt64(c echo.Context) (int64, error) {
	logger := log.GetLogger(c).Sugar()

	// IDをstringからint64に変換
	dogOwnerID, err := strconv.ParseInt(claims.ID, 10, 64)
	if err != nil {
		wrErr := errors.NewWRError(
			nil,
			"認証情報が違います。",
			errors.NewDogOwnerClientErrorEType(),
		)
		logger.Error(wrErr)
		return 0, err
	}
	return dogOwnerID, nil
}

// LogIn: dogownerの存在チェックバリデーションとJWTの更新, 署名済みjwtを返す
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - dto.AuthDogOwnerReq: authDogOwnerのリクエスト情報
//
// return:
//   - string: 検証済みのjwt
//   - error: error情報
func (ah *authHandler) LogIn(c echo.Context, ador authDTO.AuthDogOwnerReq) (string, error) {
	logger := log.GetLogger(c).Sugar()

	// EmailとPhoneNumberのバリデーション
	if wrErr := validateEmailOrPhoneNumber(ador); wrErr != nil {
		logger.Error(wrErr)
		return "", wrErr
	}

	logger.Debugf("authDogOwnerReq: %v, Type: %T", ador, ador)

	// EmailかPhoneNumberから対象のDogOwner情報の取得
	result, err := ah.ar.GetDogOwnerByCredential(c, ador)

	if err != nil {
		logger.Error(err)
		return "", err
	}

	// パスワードの確認
	err = bcrypt.CompareHashAndPassword([]byte(result.Password.String), []byte(ador.Password))

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"パスワードが間違っています",
			wrErrors.NewDogOwnerServerErrorEType())

		logger.Errorf("Password compare failure: %v", wrErr)

		return "", wrErr
	}

	// 更新用のJWT IDの生成
	jwtID, wrErr := GenerateJwtID(c)

	if wrErr != nil {
		return "", wrErr
	}

	// 取得したdogOwnerのjtw_idの更新
	wrErr = ah.ar.UpdateJwtID(c, result, jwtID)

	if wrErr != nil {
		return "", wrErr
	}

	// 作成したDogOwnerの情報をdto詰め替え
	dogOwnerDetail := doDTO.DogOwnerDTO{
		DogOwnerID: result.AuthDogOwner.DogOwnerID.Int64,
		JwtID:      jwtID,
	}

	logger.Infof("dogOwnerDetail: %v", dogOwnerDetail)

	// 署名済みのjwt token取得
	token, wrErr := GetSignedJwt(c, dogOwnerDetail)

	if wrErr != nil {
		return "", wrErr
	}

	return token, nil
}

// Revoke: Revoke機能
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *dogMW.AccountClaims: 検証済みのclaims情報
//
// return:
//   - error: error情報
func (ah *authHandler) Revoke(c echo.Context, claims *AccountClaims) error {
	// dogOwnerIDの取得
	dogOwnerID, wrErr := claims.GetDogOwnerIDAsInt64(c)
	if wrErr != nil {
		return wrErr
	}

	// 対象のdogOwnerのIDからJWT IDの削除
	if wrErr := ah.ar.DeleteJwtID(c, dogOwnerID); wrErr != nil {
		return wrErr
	}

	return nil
}

/*
Google OAuth認証
*/
// func (ah *authHandler) GoogleOAuth(c echo.Context, authorizationCode string, grantType types.GrantType) (dto.ResDogOwnerDto, error) {
// 	logger := log.GetLogger(c).Sugar()

// 	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second) // 5秒で設定
// 	defer cancel()

// 	// 各token情報の取得
// 	token, wrErr := ah.ag.GetAccessToken(c, authorizationCode, ctx)

// 	if wrErr != nil {
// 		return dto.ResDogOwnerDto{}, wrErr
// 	}

// 	// トークン元にGoogleユーザー情報の取得
// 	googleUserInfo, wrErr := ah.ag.GetGoogleUserInfo(c, token, ctx)

// 	if wrErr != nil {
// 		return dto.ResDogOwnerDto{}, wrErr
// 	}

// 	// Googleユーザー情報の確認処理
// 	if googleUserInfo == nil {
// 		wrErr := wrErrors.NewWRError(
// 			errors.New(""),
// 			"no google user information",
// 			wrErrors.NewAuthServerErrorEType(),
// 		)
// 		logger.Errorf("No google user information error: %v", wrErr)
// 		return dto.ResDogOwnerDto{}, wrErr
// 	}

// 	// ドッグオーナーのcredentialの設定と型変換
// 	dogOwnerCredential := model.DogOwnerCredential{
// 		ProviderUserID: wrUtil.NewSqlNullString(googleUserInfo.UserId),
// 		Email:          wrUtil.NewSqlNullString(googleUserInfo.Email),
// 		AuthDogOwner: model.AuthDogOwner{
// 			AccessToken:           wrUtil.NewSqlNullString(token.AccessToken),
// 			RefreshToken:          wrUtil.NewSqlNullString(token.RefreshToken),
// 			AccessTokenExpiration: wrUtil.NewCustomTime(token.Expiry),
// 			GrantType:             grantType,
// 			DogOwner: model.DogOwner{
// 				Name: wrUtil.NewSqlNullString(googleUserInfo.Email),
// 			},
// 		},
// 	}

// 	// ドッグオーナーの作成
// 	dogOC, wrErr := ah.ar.CreateOAuthDogOwner(c, &dogOwnerCredential)

// 	if wrErr != nil {
// 		return dto.ResDogOwnerDto{}, wrErr
// 	}

// 	resDogOwner := dto.ResDogOwnerDto{
// 		DogOwnerID: dogOC.AuthDogOwner.DogOwner.DogOwnerID.Int64,
// 	}

// 	return resDogOwner, nil
// }

/*
jwt処理
*/
// GetSignedJwt: 署名済みのJWT tokenの取得
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - dto.DogOwnerDTO: 作成したdogOwnerの情報
//
// return:
//  - string: 署名したtoken
//  - error: error情報
func GetSignedJwt(c echo.Context, dod doDTO.DogOwnerDTO) (string, error) {
	// 秘密鍵取得
	secretKey := configs.FetchConfigStr("jwt.os.secret.key")
	jwtExpTime := configs.FetchConfigInt("jwt.exp.time")

	// jwt token生成
	signedToken, wrErr := createToken(c, secretKey, dod, jwtExpTime)

	if wrErr != nil {
		return "", wrErr
	}

	return signedToken, wrErr
}

// createToken: 指定された秘密鍵を使用して認証用のJWTトークンを生成
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - string: secretKey   トークンの署名に使用する秘密鍵を表す文字列
//   - dto.DogOwnerDTO:  作成したdogOwnerの情報
//   - int: expTime トークンの有効期限を秒単位で指定
//
// return:
//   - string: 生成されたJWTトークンを表す文字列
//   - error: トークンの生成中に問題が発生したエラー
func createToken(c echo.Context, secretKey string, dod doDTO.DogOwnerDTO, expTime int) (string, error) {
	logger := log.GetLogger(c).Sugar()

	// JWTのペイロード
	claims := AccountClaims{
		ID:  strconv.FormatInt(dod.DogOwnerID, 10), // stringにコンバート
		JTI: dod.JwtID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expTime))), // 有効時間
		},
	}

	// token生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// tokenに署名
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"パスワードに不正な文字列が入っています。",
			wrErrors.NewDogOwnerClientErrorEType(),
		)
		logger.Error(wrErr)
		return "", err
	}

	return signedToken, nil
}

// GenerateJwtID: JwtIDの生成。引数の数だけランダムの文字列を生成
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//
// return:
//   - string: JwtID
//   - error: error情報
func GenerateJwtID(c echo.Context) (string, error) {
	logger := log.GetLogger(c).Sugar()

	// カスタムエラー処理
	handleError := func(err error) error {
		wrErr := wrErrors.NewWRError(
			err,
			"JwtID生成に失敗しました",
			wrErrors.NewDogOwnerServerErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	// UUIDを生成
	return util.UUIDGenerator(handleError)
}

// validateEmailOrPhoneNumber: EmailかPhoneNumberの識別バリデーション。パスワード認証は、EmailかPhoneNumberで登録するため
//
// args:
//   - dto.DogOwnerReq: DogOwnerのRequest
//
// return:
//   - error: err情報
func validateEmailOrPhoneNumber(doReq authDTO.AuthDogOwnerReq) error {
	// 両方が空の場合はエラー
	if doReq.Email == "" && doReq.PhoneNumber == "" {
		wrErr := wrErrors.NewWRError(
			nil,
			"Emailと電話番号のどちらも空です",
			wrErrors.NewDogOwnerClientErrorEType(),
		)
		return wrErr
	}

	// 両方に値が入っている場合もエラー
	if doReq.Email != "" && doReq.PhoneNumber != "" {
		wrErr := wrErrors.NewWRError(
			nil,
			"Emailと電話番号のどちらも値が入っています",
			wrErrors.NewDogOwnerClientErrorEType(),
		)
		return wrErr
	}

	// どちらか片方だけが入力されている場合は正常
	return nil
}
