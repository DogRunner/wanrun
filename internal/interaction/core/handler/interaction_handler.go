package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dogrun/facade"
	"github.com/wanrun-develop/wanrun/internal/interaction/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/interaction/core/dto"
	"github.com/wanrun-develop/wanrun/internal/wrcontext"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IBookmarkHandler interface {
	AddBookmark(echo.Context, dto.BookmarkAddReq) ([]int64, error)
	DeleteBookmark(echo.Context, dto.BookmarkDeleteReq) error
}

type bookmarkHandler struct {
	r  repository.IBookmarkRepository
	df facade.IDogrunFacade
}

func NewBookmarkHandler(br repository.IBookmarkRepository, df facade.IDogrunFacade) IBookmarkHandler {
	return &bookmarkHandler{br, df}
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
func (h *bookmarkHandler) AddBookmark(c echo.Context, reqBody dto.BookmarkAddReq) ([]int64, error) {

	logger := log.GetLogger(c).Sugar()
	logger.Info("dogrunのお気に入り登録. dogrunID: ", reqBody.DogrunIDs)
	err := h.df.CheckDogrunExistByIDs(c, reqBody.DogrunIDs)
	if err != nil {
		return nil, err
	}
	logger.Info("dogrunの存在チェック済み")

	//ログインユーザーIDの取得
	userID, err := wrcontext.GetLoginUserID(c)
	if err != nil {
		return nil, err
	}

	bookmarkIDs := []int64{}

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
		bookmarkIDs = append(bookmarkIDs, bookmarkId)
	}

	return bookmarkIDs, nil
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
func (h *bookmarkHandler) DeleteBookmark(c echo.Context, reqBody dto.BookmarkDeleteReq) error {
	logger := log.GetLogger(c).Sugar()
	logger.Info("dogrunのお気に入り削除. dogrunID: ", reqBody.DogrunIDs)

	//ログインユーザーIDの取得
	userID, err := wrcontext.GetLoginUserID(c)
	if err != nil {
		return err
	}

	//削除処理
	if err := h.r.DeleteBookmark(c, reqBody.DogrunIDs, userID); err != nil {
		return err
	}
	return nil
}
