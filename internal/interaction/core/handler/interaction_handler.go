package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dogrun/core/handler"
	"github.com/wanrun-develop/wanrun/internal/interaction/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/interaction/core/dto"
	"github.com/wanrun-develop/wanrun/internal/wrcontext"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IBookmarkHandler interface {
	AddBookmark(echo.Context, dto.AddBookmark) ([]int64, error)
	DeleteBookmark(echo.Context, dto.DeleteBookmark) error
}

type bookmarkHandler struct {
	r   repository.IBookmarkRepository
	drh handler.IDogrunHandler
}

func NewBookmarkHandler(br repository.IBookmarkRepository, dh handler.IDogrunHandler) IBookmarkHandler {
	return &bookmarkHandler{br, dh}
}

// AddBookmark: ブックマークへのdogrunの追加
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.AddBookmark:	リクエストボディ
//
// return:
//   - int:	int
//   - error:	エラー
func (h *bookmarkHandler) AddBookmark(c echo.Context, reqBody dto.AddBookmark) ([]int64, error) {

	logger := log.GetLogger(c).Sugar()
	logger.Info("dogrunのお気に入り登録. dogrunID: ", reqBody.DogrunIDs)
	err := h.drh.CheckDogrunExistByIds(c, reqBody.DogrunIDs)
	if err != nil {
		return nil, err
	}
	logger.Info("dogrunの存在チェック済み")

	//ログインユーザーIDの取得
	userID, err := wrcontext.GetLoginUserId(c)
	if err != nil {
		return nil, err
	}

	bookmarkIds := []int64{}

	//ひとつずつ、すでにブックマーク済みかチェック
	for _, dogrunID := range reqBody.DogrunIDs {
		bookmark, err := h.r.FindDogrunBookmark(c, dogrunID, userID)
		if err != nil {
			return nil, err
		}
		if bookmark.IsNotEmpty() {
			err = errors.NewWRError(nil, fmt.Sprintf("ドッグランID:%dはすでにブックマークに登録されています。", dogrunID), errors.NewInteractionClientErrorEType())
			logger.Error("ブックマーク既存チェックでバリデーションエラー", err)
			return nil, err
		}
		//ブックマークに登録
		bookmarkId, err := h.r.AddBookmark(c, dogrunID, userID)
		if err != nil {
			return nil, err
		}
		bookmarkIds = append(bookmarkIds, bookmarkId)
	}

	return bookmarkIds, nil
}

// DeleteBookmark: ブックマークへのdogrunの追加
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.DeleteBookmark:	リクエストボディ
//
// return:
//   - int:	int
//   - error:	エラー
func (h *bookmarkHandler) DeleteBookmark(c echo.Context, reqBody dto.DeleteBookmark) error {
	logger := log.GetLogger(c).Sugar()
	logger.Info("dogrunのお気に入り削除. dogrunID: ", reqBody.DogrunIDs)

	//ログインユーザーIDの取得
	userID, err := wrcontext.GetLoginUserId(c)
	if err != nil {
		return err
	}

	//削除処理
	if err := h.r.DeleteBookmark(c, reqBody.DogrunIDs, userID); err != nil {
		return err
	}
	return nil
}
