package facade

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dog/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/wrcontext"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IDogFacade interface {
	CheckDogownerValid(echo.Context, []int64) error
}

type dogFacade struct {
	dr repository.IDogRepository
}

func NewDogFacade(drr repository.IDogRepository) IDogFacade {
	return &dogFacade{drr}
}

// CheckDogownerValid: dogのdogownerが正しいかチェック
// ログインユーザーのdogであるかをチェック
//
// args:
//   - echo.Context:	コンテキスト
//   - []]nt64:	dogIDs チェック対象のdogIDs
//
// return:
// error:	エラー
func (f dogFacade) CheckDogownerValid(c echo.Context, paramDogIDs []int64) error {
	logger := log.GetLogger(c).Sugar()

	// ログインユーザーIDの取得
	userID, err := wrcontext.GetLoginUserID(c)
	if err != nil {
		return err
	}
	//ユーザーIDを条件にdog取得
	dogsResults, err := f.dr.GetDogByDogOwnerID(c, userID)
	if err != nil {
		return err
	}
	//チェック用一時map
	dogMap := make(map[int64]struct{})
	for _, dog := range dogsResults {
		dogMap[dog.DogID.Int64] = struct{}{}
	}

	for _, paramDogID := range paramDogIDs {
		if _, exists := dogMap[paramDogID]; !exists {
			err := errors.NewWRError(nil, fmt.Sprintf("指定されたドッグID:%dはあなたのペットではありません", paramDogID), errors.NewInteractionClientErrorEType())
			logger.Error(err, "不正なdogIdの指定")
			return err
		}
	}

	return nil
}
