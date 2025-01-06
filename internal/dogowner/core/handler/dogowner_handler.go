package handler

import (
	"github.com/labstack/echo/v4"
	authRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	authDTO "github.com/wanrun-develop/wanrun/internal/auth/core/dto"
	authHandler "github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	doDTO "github.com/wanrun-develop/wanrun/internal/dogowner/core/dto"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/internal/transaction"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	wrUtil "github.com/wanrun-develop/wanrun/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IDogownerHandler interface {
	DogownerSignUp(c echo.Context, doReq doDTO.DogownerReq) (string, error)
}

type dogownerHandler struct {
	dosr dogownerRepository.IDogownerScopeRepository
	tm   transaction.ITransactionManager
	asr  authRepository.IAuthScopeRepository
	dor  dogownerRepository.IDogownerRepository
	ar   authRepository.IAuthRepository
}

func NewDogownerHandler(
	dosr dogownerRepository.IDogownerScopeRepository,
	tm transaction.ITransactionManager,
	asr authRepository.IAuthScopeRepository,
	dor dogownerRepository.IDogownerRepository,
	ar authRepository.IAuthRepository,
) IDogownerHandler {
	return &dogownerHandler{
		dosr: dosr,
		tm:   tm,
		asr:  asr,
		dor:  dor,
		ar:   ar,
	}
}

// DogownerSignUp: dogownerの登録し、検証済みのJWTを返す
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//   - doDTO.DogownerReq: dogownerに対するリクエスト情報
//
// return:
//   - string
//   - error: error情報
func (doh *dogownerHandler) DogownerSignUp(c echo.Context, doReq doDTO.DogownerReq) (string, error) {
	logger := log.GetLogger(c).Sugar()

	// パスワードのハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(doReq.Password), bcrypt.DefaultCost) // 一旦costをデフォルト値

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"パスワードに不正な文字列が入っています。",
			wrErrors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return "", wrErr
	}

	// EmailとPhoneNumberのバリデーション
	if wrErr := validateEmailOrPhoneNumber(doReq); wrErr != nil {
		logger.Error(wrErr)
		return "", wrErr
	}

	// JWT IDの生成
	jwtID, wrErr := authHandler.GenerateJwtID(c)

	if wrErr != nil {
		return "", wrErr
	}

	// requestからDogownerの構造体に詰め替え
	dogownerCredential := model.DogownerCredential{
		Email:       wrUtil.NewSqlNullString(doReq.Email),
		PhoneNumber: wrUtil.NewSqlNullString(doReq.PhoneNumber),
		Password:    wrUtil.NewSqlNullString(string(hash)),
		GrantType:   wrUtil.NewSqlNullString(model.PASSWORD_GRANT_TYPE), // Password認証
		AuthDogowner: model.AuthDogowner{
			JwtID: wrUtil.NewSqlNullString(jwtID),
			Dogowner: model.Dogowner{
				Name: wrUtil.NewSqlNullString(doReq.DogownerName),
			},
		},
	}

	logger.Debugf("dogownerCredential %v, Type: %T", dogownerCredential, dogownerCredential)

	ctx := c.Request().Context()

	// Emailの重複チェック
	if wrErr := doh.ar.CheckDuplicate(c, model.EmailField, dogownerCredential.Email); wrErr != nil {
		return "", wrErr
	}

	// PhoneNumberの重複チェック
	if wrErr := doh.ar.CheckDuplicate(c, model.PhoneNumberField, dogownerCredential.PhoneNumber); wrErr != nil {
		return "", wrErr
	}

	// dogownerの作成する1トランザクション
	if err := doh.tm.DoInTransaction(c, ctx, func(tx *gorm.DB) error {

		// Dogownerを作成
		if wrErr := doh.dosr.CreateDogowner(tx, c, &dogownerCredential); wrErr != nil {
			return wrErr
		}

		// AuthDogownerを作成
		if wrErr := doh.asr.CreateAuthDogowner(tx, c, &dogownerCredential); wrErr != nil {
			return wrErr
		}

		// DogownerのCredentialを作成
		if wrErr := doh.asr.CreateDogownerCredential(tx, c, &dogownerCredential); wrErr != nil {
			return wrErr
		}

		// 正常に完了
		return nil

	}); err != nil {
		logger.Error("Transaction failed:", err)
		return "", err
	}

	// 正常に終了
	logger.Infof("Successfully created SignUp Dogowner: %v", dogownerCredential)

	// 作成したDogownerの情報をdto詰め替え
	dogownerDetail := authDTO.UserAuthInfoDTO{
		UserID: dogownerCredential.AuthDogowner.DogownerID.Int64,
		JwtID:  dogownerCredential.AuthDogowner.JwtID.String,
		RoleID: authHandler.DOGOWNER_ROLE,
	}

	logger.Infof("dogownerDetail: %v", dogownerDetail)

	// 署名済みのjwt token取得
	token, wrErr := authHandler.GetSignedJwt(c, dogownerDetail)

	if wrErr != nil {
		return "", wrErr
	}

	return token, nil
}

// validateEmailOrPhoneNumber: EmailかPhoneNumberの識別バリデーション。パスワード認証は、EmailかPhoneNumberで登録するため
//
// args:
//   - dto.DogownerReq: DogownerのRequest
//
// return:
//   - error: err情報
func validateEmailOrPhoneNumber(doReq doDTO.DogownerReq) error {
	// 両方が空の場合はエラー
	if doReq.Email == "" && doReq.PhoneNumber == "" {
		wrErr := wrErrors.NewWRError(
			nil,
			"Emailと電話番号のどちらも空です",
			wrErrors.NewDogownerClientErrorEType(),
		)
		return wrErr
	}

	// 両方に値が入っている場合もエラー
	if doReq.Email != "" && doReq.PhoneNumber != "" {
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
