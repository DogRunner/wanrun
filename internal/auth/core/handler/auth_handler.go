package handler

import (
	_ "context"
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/configs"
	"github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/auth/core/dto"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	wrUtil "github.com/wanrun-develop/wanrun/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

type IAuthHandler interface {
	CreateDogOwner(c echo.Context, ador dto.AuthDogOwnerReq) (dto.DogOwnerDTO, error)
	GetSignedJwt(c echo.Context, dod dto.DogOwnerDTO) (string, error)
	FetchDogOwnerInfo(c echo.Context, ador dto.AuthDogOwnerReq) (dto.DogOwnerDTO, error)
	// LogOut() error
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

// CreateDogOwner: DogOwnerの作成
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - dto.AuthDogOwnerReq: authDogOwnerのリクエスト情報
//
// return:
//   - dto.DogOwnerDOT: 作成したdogOwner情報
//   - error: error情報
func (ah *authHandler) CreateDogOwner(c echo.Context, ador dto.AuthDogOwnerReq) (dto.DogOwnerDTO, error) {
	logger := log.GetLogger(c).Sugar()

	// パスワードのハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(ador.Password), bcrypt.DefaultCost) // 一旦costをデフォルト値

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"パスワードに不正な文字列が入っております。",
			wrErrors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return dto.DogOwnerDTO{}, wrErr
	}

	// EmailとPhoneNumberのバリデーション
	if wrErr := validateEmailOrPhoneNumber(ador); wrErr != nil {
		logger.Error(wrErr)
		return dto.DogOwnerDTO{}, wrErr
	}

	// JWT IDの生成
	jwtID, wrErr := createJwtID(c, 15)

	if wrErr != nil {
		return dto.DogOwnerDTO{}, wrErr
	}

	// requestからDogOwnerの構造体に詰め替え
	dogOwnerCredential := model.DogOwnerCredential{
		Email:       wrUtil.NewSqlNullString(ador.Email),
		PhoneNumber: wrUtil.NewSqlNullString(ador.PhoneNumber),
		Password:    wrUtil.NewSqlNullString(string(hash)),
		GrantType:   wrUtil.NewSqlNullString(model.PASSWORD_GRANT_TYPE), // Password認証
		AuthDogOwner: model.AuthDogOwner{
			JwtID: wrUtil.NewSqlNullString(jwtID),
			DogOwner: model.DogOwner{
				Name: wrUtil.NewSqlNullString(ador.DogOwnerName),
			},
		},
	}

	logger.Debugf("dogOwnerCredential %v, Type: %T", dogOwnerCredential, dogOwnerCredential)

	// ドッグオーナー作成
	result, wrErr := ah.ar.CreateDogOwner(c, &dogOwnerCredential)

	if wrErr != nil {
		return dto.DogOwnerDTO{}, wrErr
	}

	// 作成したDogOwnerの情報をdto詰め替え
	dogOwnerDetail := dto.DogOwnerDTO{
		DogOwnerID: result.AuthDogOwner.DogOwnerID.Int64,
		JwtID:      result.AuthDogOwner.JwtID.String,
	}

	logger.Infof("dogOwnerDetail: %v", dogOwnerDetail)

	return dogOwnerDetail, nil
}

// FetchDogOwnerInfo: リクエストのバリデーションやdogOwnerの情報取得
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - dto.AuthDogOwnerReq: authDogOwnerのリクエスト情報
//
// return:
//   - dto.DogOwnerDOT: 作成したdogOwner情報
//   - error: error情報
func (ah *authHandler) FetchDogOwnerInfo(c echo.Context, ador dto.AuthDogOwnerReq) (dto.DogOwnerDTO, error) {
	logger := log.GetLogger(c).Sugar()

	// EmailとPhoneNumberのバリデーション
	if wrErr := validateEmailOrPhoneNumber(ador); wrErr != nil {
		logger.Error(wrErr)
		return dto.DogOwnerDTO{}, wrErr
	}

	logger.Debugf("authDogOwnerReq: %v, Type: %T", ador, ador)

	// EmailかPhoneNumberから対象のDogOwner情報の取得
	result, err := ah.ar.GetDogOwnerByCredential(c, ador)

	if err != nil {
		logger.Error(err)
		return dto.DogOwnerDTO{}, err
	}

	// パスワードの確認
	err = bcrypt.CompareHashAndPassword([]byte(result.Password.String), []byte(ador.Password))

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"パスワードが間違っています",
			wrErrors.NewDogownerServerErrorEType())

		logger.Errorf("Password compare failure: %v", wrErr)

		return dto.DogOwnerDTO{}, wrErr
	}

	// 更新用のJWT IDの生成
	jwtID, wrErr := createJwtID(c, 15)

	if wrErr != nil {
		return dto.DogOwnerDTO{}, wrErr
	}

	// 取得したdogOwnerのjtw_idの更新
	wrErr = ah.ar.UpdateJwtID(c, result, jwtID)

	if wrErr != nil {
		return dto.DogOwnerDTO{}, wrErr
	}

	// 作成したDogOwnerの情報をdto詰め替え
	dogOwnerDetail := dto.DogOwnerDTO{
		DogOwnerID: result.AuthDogOwner.DogOwnerID.Int64,
		JwtID:      jwtID,
	}

	logger.Infof("dogOwnerDetail: %v", dogOwnerDetail)

	return dogOwnerDetail, nil
}

// Logout
func (ah *authHandler) LogOut() error { return nil }

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

// validateEmailOrPhoneNumber: EmailかPhoneNumberの識別バリデーション。パスワード認証は、EmailかPhoneNumberで登録するため
//
// args:
//   - dto.ReqAuthDogOwnerDto: Response用のAuthDogOwnerのDTO
//
// return:
//   - error: err情報
func validateEmailOrPhoneNumber(ador dto.AuthDogOwnerReq) error {
	// 両方が空の場合はエラー
	if ador.Email == "" && ador.PhoneNumber == "" {
		wrErr := wrErrors.NewWRError(
			nil,
			"Emailと電話番号のどちらも空です",
			wrErrors.NewDogownerClientErrorEType(),
		)
		return wrErr
	}

	// 両方に値が入っている場合もエラー
	if ador.Email != "" && ador.PhoneNumber != "" {
		wrErr := wrErrors.NewWRError(
			nil,
			"Emailと電話番号のどちらも値が入っています",
			wrErrors.NewDogownerClientErrorEType(),
		)
		return wrErr
	}

	// どちらか片方だけが入力されている場合は正常
	return nil
}

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

func (ah *authHandler) GetSignedJwt(c echo.Context, dod dto.DogOwnerDTO) (string, error) {
	// 秘密鍵取得
	secretKey := configs.FetchCondigStr("jwt.os.secret.key")
	jwtExpTime := configs.FetchCondigInt("jwt.exp.time")

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
func createToken(c echo.Context, secretKey string, dod dto.DogOwnerDTO, expTime int) (string, error) {
	logger := log.GetLogger(c).Sugar()

	// JWTのペイロード
	claims := &dto.AccountClaims{
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
			wrErrors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return "", err
	}

	return signedToken, nil
}

// createJwtID: JWT IDの生成。引数の数だけランダムの文字列を生成
//
// args:
//   - int: length 生成したい数
//
// return:
//   - string: ランダム文字列
//   - error: error情報
func createJwtID(c echo.Context, length int) (string, error) {
	logger := log.GetLogger(c).Sugar()

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"JWT ID生成に失敗しました",
			wrErrors.NewDogownerServerErrorEType(),
		)
		logger.Error(wrErr)
		return "", wrErr
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}
