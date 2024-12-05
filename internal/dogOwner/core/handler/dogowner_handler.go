package handler

import (
	"github.com/labstack/echo/v4"
	authRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	authScopeRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/scopeRepository"
	dogOwnerRepository "github.com/wanrun-develop/wanrun/internal/dogOwner/adapters/repository"
	dogOwnerScopeRepository "github.com/wanrun-develop/wanrun/internal/dogOwner/adapters/scopeRepository"
	"github.com/wanrun-develop/wanrun/internal/dogOwner/core/dto"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/internal/transaction"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
	wrUtil "github.com/wanrun-develop/wanrun/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IDogOwnerHandler interface {
	DogOwnerSignUp(c echo.Context, doReq dto.DogOwnerReq) (dto.DogOwnerDTO, error)
}

type dogOwnerHandler struct {
	dosr dogOwnerScopeRepository.IDogOwnerScopeRepository
	tm   transaction.ITransactionManager
	asr  authScopeRepository.IAuthScopeRepository
	dor  dogOwnerRepository.IDogOwnerRepository
	ar   authRepository.IAuthRepository
}

func NewDogOwnerHandler(
	dosr dogOwnerScopeRepository.IDogOwnerScopeRepository,
	tm transaction.ITransactionManager,
	asr authScopeRepository.IAuthScopeRepository,
	dor dogOwnerRepository.IDogOwnerRepository,
	ar authRepository.IAuthRepository,
) IDogOwnerHandler {
	return &dogOwnerHandler{
		dosr: dosr,
		tm:   tm,
		asr:  asr,
		dor:  dor,
		ar:   ar,
	}
}

// DogOwnerSignUp: dogOwnerの登録処理
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//   - dto.DogOwnerReq: dogOwnerに対するリクエスト情報
//
// return:
//   - dto.dogOwnerDTO: dogOwnerのレスポンス情報
//   - error: error情報
func (doh *dogOwnerHandler) DogOwnerSignUp(c echo.Context, doReq dto.DogOwnerReq) (dto.DogOwnerDTO, error) {
	logger := log.GetLogger(c).Sugar()

	// パスワードのハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(doReq.Password), bcrypt.DefaultCost) // 一旦costをデフォルト値

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"パスワードに不正な文字列が入っております。",
			wrErrors.NewDogOwnerClientErrorEType(),
		)
		logger.Error(wrErr)
		return dto.DogOwnerDTO{}, wrErr
	}

	// EmailとPhoneNumberのバリデーション
	if wrErr := validateEmailOrPhoneNumber(doReq); wrErr != nil {
		logger.Error(wrErr)
		return dto.DogOwnerDTO{}, wrErr
	}

	// JWT IDの生成
	jwtID, wrErr := generateJwtID(c, 15)

	if wrErr != nil {
		return dto.DogOwnerDTO{}, wrErr
	}

	// requestからDogOwnerの構造体に詰め替え
	dogOwnerCredential := model.DogOwnerCredential{
		Email:       wrUtil.NewSqlNullString(doReq.Email),
		PhoneNumber: wrUtil.NewSqlNullString(doReq.PhoneNumber),
		Password:    wrUtil.NewSqlNullString(string(hash)),
		GrantType:   wrUtil.NewSqlNullString(model.PASSWORD_GRANT_TYPE), // Password認証
		AuthDogOwner: model.AuthDogOwner{
			JwtID: wrUtil.NewSqlNullString(jwtID),
			DogOwner: model.DogOwner{
				Name: wrUtil.NewSqlNullString(doReq.DogOwnerName),
			},
		},
	}

	logger.Debugf("dogOwnerCredential %v, Type: %T", dogOwnerCredential, dogOwnerCredential)

	ctx := c.Request().Context()

	// Emailの重複チェック
	if wrErr := doh.ar.CheckDuplicate(c, model.EmailField, dogOwnerCredential.Email); wrErr != nil {
		return dto.DogOwnerDTO{}, wrErr
	}

	// PhoneNumberの重複チェック
	if wrErr := doh.ar.CheckDuplicate(c, model.PhoneNumberField, dogOwnerCredential.PhoneNumber); wrErr != nil {
		return dto.DogOwnerDTO{}, wrErr
	}

	// dogOwnerの作成する1トランザクション
	if err := doh.tm.DoInTransaction(c, ctx, func(tx *gorm.DB) error {

		// DogOwnerを作成
		if wrErr := doh.dosr.CreateDogOwner(tx, c, &dogOwnerCredential); wrErr != nil {
			return wrErr
		}

		// AuthDogOwnerを作成
		if wrErr := doh.asr.CreateAuthDogOwner(tx, c, &dogOwnerCredential); wrErr != nil {
			return wrErr
		}

		// DogOwnerのCredentialを作成
		if wrErr := doh.asr.CreateDogOwnerCredential(tx, c, &dogOwnerCredential); wrErr != nil {
			return wrErr
		}

		// 正常に完了
		return nil

	}); err != nil {
		logger.Error("Transaction failed:", err)
		return dto.DogOwnerDTO{}, err
	}

	// 正常に終了
	logger.Infof("Successfully created SignUp DogOwner: %v", dogOwnerCredential)

	// 作成したDogOwnerの情報をdto詰め替え
	dogOwnerDetail := dto.DogOwnerDTO{
		DogOwnerID: dogOwnerCredential.AuthDogOwner.DogOwnerID.Int64,
		JwtID:      dogOwnerCredential.AuthDogOwner.JwtID.String,
	}

	logger.Infof("dogOwnerDetail: %v", dogOwnerDetail)

	return dogOwnerDetail, nil
}

// validateEmailOrPhoneNumber: EmailかPhoneNumberの識別バリデーション。パスワード認証は、EmailかPhoneNumberで登録するため
//
// args:
//   - dto.DogOwnerReq: DogOwnerのRequest
//
// return:
//   - error: err情報
func validateEmailOrPhoneNumber(doReq dto.DogOwnerReq) error {
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

// generateJwtID: JwtIDの生成。引数の数だけランダムの文字列を生成
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - int: 生成されるIDの長さを指定
//
// return:
//   - string: JwtID
//   - error: error情報
func generateJwtID(c echo.Context, l int) (string, error) {
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
	return util.UUIDGenerator(l, handleError)
}
