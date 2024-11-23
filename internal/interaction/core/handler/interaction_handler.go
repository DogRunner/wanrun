package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dogrun/core/handler"
	"github.com/wanrun-develop/wanrun/internal/interaction/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/interaction/core/dto"
	"github.com/wanrun-develop/wanrun/internal/wrcontext"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

type IBookmarkHandler interface {
	AddBookmark(echo.Context, dto.AddBookmark) (int64, error)
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
//   - dto.AddFavorite:	リクエストボディ
//
// return:
//   - int:	int
//   - error:	エラー
func (h *bookmarkHandler) AddBookmark(c echo.Context, reqBody dto.AddBookmark) (int64, error) {

	logger := log.GetLogger(c).Sugar()
	logger.Info("dogrunのお気に入り登録. dogrunID: ", reqBody.DogrunID)
	if err := h.drh.CheckDogrunExistById(c, reqBody.DogrunID); err != nil {
		return 0, err
	}
	logger.Info("dogrunの存在チェック済み")

	//ログインユーザーIDの取得
	userID, err := wrcontext.GetLoginUserId(c)
	if err != nil {
		return 0, nil
	}

	//すでにブックマーク済みかチェック
	bookmark, err := h.r.FindDogrunBookmark(c, reqBody.DogrunID, userID)
	if err != nil {
		return 0, nil
	}
	if bookmark.IsNotEmpty() {
		err = errors.NewWRError(nil, "すでにブックマークに登録されています。", errors.NewInteractionClientErrorEType())
		logger.Error("ブックマーク既存チェックでバリデーションエラー", err)
		return 0, err
	}

	//ブックマークに登録
	bookmarkId, err := h.r.AddBookmark(c, reqBody.DogrunID, userID)
	if err != nil {
		return 0, err
	}

	return bookmarkId, nil
}
