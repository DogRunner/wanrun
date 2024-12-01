package usecase

import (
	"github.com/labstack/echo/v4"
	authRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	dogOwnerRepository "github.com/wanrun-develop/wanrun/internal/dogOwner/adapters/repository"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/internal/transaction"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IDogOwnerUsecase interface {
	SignUp(c echo.Context, doc *model.DogOwnerCredential) (*model.DogOwnerCredential, error)
}

type dogOwnerUsecase struct {
	ar  authRepository.IAuthRepository
	dor dogOwnerRepository.IDogOwnerRepository
	tm  transaction.ITransactionManager
}

func NewDogOwnerUsecase(
	ar authRepository.IAuthRepository,
	dor dogOwnerRepository.IDogOwnerRepository,
	tm transaction.ITransactionManager,
) IDogOwnerUsecase {
	return &dogOwnerUsecase{
		ar:  ar,
		dor: dor,
		tm:  tm,
	}
}

// SignUp: dogOwnerの登録処理
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用されます。
//   - *model.DogOwnerCredential: dogOwnerの情報
//
// return:
//   - *model.DogOwnerCredential: dogOwnerの情報
//   - error: error情報
func (dou *dogOwnerUsecase) SignUp(c echo.Context, doc *model.DogOwnerCredential) (*model.DogOwnerCredential, error) {
	logger := log.GetLogger(c).Sugar()

	ctx := c.Request().Context()

	// dogOwnerの作成する1トランザクション
	err := dou.tm.DoInTransaction(c, ctx, func(tx *gorm.DB) error {

		// Emailの重複チェック
		if wrErr := dou.ar.CheckDuplicate(c, model.EmailField, doc.Email); wrErr != nil {
			return wrErr
		}

		// PhoneNumberの重複チェック
		if wrErr := dou.ar.CheckDuplicate(c, model.PhoneNumberField, doc.PhoneNumber); wrErr != nil {
			return wrErr
		}

		// DogOwnerを作成
		if wrErr := dou.dor.CreateDogOwner(tx, c, doc); wrErr != nil {
			return wrErr
		}

		// AuthDogOwnerとDogOwnerのCredential作成
		if wrErr := dou.ar.CreateAuthDogOwnerAndCredential(tx, c, doc); wrErr != nil {
			return wrErr
		}

		// 正常に完了
		return nil

	})

	if err != nil {
		logger.Error("Transaction failed:", err)
		return nil, err
	}

	// 正常に終了
	logger.Infof("Successfully created SignUp DogOwner: %v", doc)
	return doc, nil
}
