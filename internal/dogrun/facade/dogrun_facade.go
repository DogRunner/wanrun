package facade

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dogrun/adapters/repository"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IDogrunFacade interface {
	CheckDogrunExistByIDs(echo.Context, []int64) error
}

type dogrunFacade struct {
	drr repository.IDogrunRepository
}

func NewDogrunFacade(drr repository.IDogrunRepository) IDogrunFacade {
	return &dogrunFacade{drr}
}

// CheckDogrunExistByIds: ドッグランの存在チェック
// args:
//   - echo.Context:	コンテキスト
//   - []int64:	ドッグランIDs
//
// return:
//   - error:	エラー
func (h *dogrunFacade) CheckDogrunExistByIDs(c echo.Context, dogrunIDs []int64) error {
	logger := log.GetLogger(c).Sugar()

	dogrunResults, err := h.drr.FindDogrunByIDs(dogrunIDs)
	if err != nil {
		err = errors.NewWRError(err, "dogrun存在チェックでエラー", errors.NewDogrunClientErrorEType())
		return err
	}

	if len(dogrunResults) != len(dogrunIDs) {
		// select結果をマップにする
		existDogrunsMap := make(map[int64]struct{})
		for _, dogrun := range dogrunResults {
			existDogrunsMap[dogrun.DogrunID.Int64] = struct{}{}
		}
		// 存在しなかったdogrunを抽出
		notExistsIDs := []int64{}
		for _, targetID := range dogrunIDs {
			if _, exists := existDogrunsMap[targetID]; !exists {
				notExistsIDs = append(notExistsIDs, targetID)
			}
		}
		err = errors.NewWRError(nil, fmt.Sprintf("指定されたドッグランID:%dが存在しません", notExistsIDs), errors.NewDogrunClientErrorEType())
		logger.Error("不正なdogrun idの指定", err)
		return err
	}
	return nil
}
