package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	dogFacade "github.com/wanrun-develop/wanrun/internal/dog/facade"
	dogrunFacade "github.com/wanrun-develop/wanrun/internal/dogrun/facade"
	"github.com/wanrun-develop/wanrun/internal/interaction/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/interaction/core/dto"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/internal/wrcontext"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
)

type IBookmarkHandler interface {
	AddBookmark(echo.Context, dto.BookmarkAddReq) ([]int64, error)
	DeleteBookmark(echo.Context, dto.BookmarkDeleteReq) error
}

type bookmarkHandler struct {
	r   repository.IBookmarkRepository
	drf dogrunFacade.IDogrunFacade
}

func NewBookmarkHandler(br repository.IBookmarkRepository, drf dogrunFacade.IDogrunFacade) IBookmarkHandler {
	return &bookmarkHandler{br, drf}
}

// AddBookmark: ブックマークへのdogrunの追加
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.BookmarkAddReq:	リクエストボディ
//
// return:
//   - int:	int
//   - error:	エラー
func (h *bookmarkHandler) AddBookmark(c echo.Context, reqBody dto.BookmarkAddReq) ([]int64, error) {

	logger := log.GetLogger(c).Sugar()
	logger.Info("dogrunのお気に入り登録. dogrunID: ", reqBody.DogrunIDs)
	err := h.drf.CheckDogrunExistByIDs(c, reqBody.DogrunIDs)
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

// DeleteBookmark: ブックマークへのdogrunの削除
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.BookmarkDeleteReq:	リクエストボディ
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

type ICheckInOutHandler interface {
	CheckinDogrun(echo.Context, dto.CheckinReq) error
	CheckoutDogrun(echo.Context, dto.CheckoutReq) error
	GetTodayCheckins(c echo.Context) ([]dto.CheckinsRes, error)
}

type checkInOutHandler struct {
	r   repository.ICheckInOutRepository
	drf dogrunFacade.IDogrunFacade
	df  dogFacade.IDogFacade
}

func NewCheckInOutHandler(br repository.ICheckInOutRepository, drf dogrunFacade.IDogrunFacade, df dogFacade.IDogFacade) ICheckInOutHandler {
	return &checkInOutHandler{br, drf, df}
}

// CheckinDogrun: ドッグランにチェックインする
// すでに一度チェックイン済みなら、re_checkin_atのみの更新
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.CheckinReq:	リクエストボディ
//
// return:
//   - error:	エラー
func (h checkInOutHandler) CheckinDogrun(c echo.Context, reqBody dto.CheckinReq) error {
	logger := log.GetLogger(c).Sugar()

	//dogrun存在チェック
	dogrunID := reqBody.DogrunID
	if err := h.drf.CheckDogrunExistByIDs(c, []int64{dogrunID}); err != nil {
		return err
	}

	//dogのdogownerチェック
	checkinDogIDs := reqBody.DogIDs
	if err := h.df.CheckDogownerValid(c, checkinDogIDs); err != nil {
		return err
	}

	saveCheckins := []model.DogrunCheckin{}
	for _, dogID := range checkinDogIDs {
		checkinResult, err := h.r.FindTodayDogrunCheckin(c, dogrunID, dogID)
		if err != nil {
			return err
		}
		if checkinResult.IsEmpty() {
			logger.Info("今日の新規チェックイン")
			checkinResult.DogrunID = util.NewSqlNullInt64(dogrunID)
			checkinResult.DogID = util.NewSqlNullInt64(dogID)
		}
		saveCheckins = append(saveCheckins, checkinResult)
	}

	//保存
	_, err := h.r.SaveDogrunCheckins(c, saveCheckins)
	if err != nil {
		return err
	}

	return nil
}

// CheckoutDogrun: ドッグランにチェックアウトする
// すでに一度チェックアウト済みならre_checkout_atのみの更新
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.CheckoutReq:	リクエストボディ
//
// return:
//   - error:	エラー
func (h checkInOutHandler) CheckoutDogrun(c echo.Context, reqBody dto.CheckoutReq) error {
	logger := log.GetLogger(c).Sugar()

	//dogrun存在チェック
	dogrunID := reqBody.DogrunID
	if err := h.drf.CheckDogrunExistByIDs(c, []int64{dogrunID}); err != nil {
		return err
	}

	//dogのdogownerチェック
	checkoutDogIDs := reqBody.DogIDs
	if err := h.df.CheckDogownerValid(c, checkoutDogIDs); err != nil {
		return err
	}

	saveCheckouts := []model.DogrunCheckout{}
	for _, dogID := range checkoutDogIDs {
		//入場しているかチェック
		checkinResult, err := h.r.FindTodayDogrunCheckin(c, dogrunID, dogID)
		if err != nil {
			return err
		}
		if checkinResult.IsEmpty() {
			err := errors.NewWRError(nil, fmt.Sprintf("今日ドッグID:%dはドッグラン%dに入場していません", dogID, dogID), errors.NewInteractionClientErrorEType())
			logger.Error(err, "入場していないドッグランへの退場リクエストのためエラー")
			return err
		}
		//すでに一度チェックアウト（退場）済みかをチェック
		checkoutResult, err := h.r.FindTodayDogrunCheckout(c, dogrunID, dogID)
		if err != nil {
			return err
		}
		if checkoutResult.IsEmpty() {
			logger.Info("今日の新規チェックアウト")
			checkoutResult.DogrunID = util.NewSqlNullInt64(dogrunID)
			checkoutResult.DogID = util.NewSqlNullInt64(dogID)
		}
		saveCheckouts = append(saveCheckouts, checkoutResult)
	}

	//保存
	_, err := h.r.SaveDogrunCheckouts(c, saveCheckouts)
	if err != nil {
		return err
	}

	return nil
}

// GetTodayCheckins: すべての所有dogの今日のチェックイン履歴の取得
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
// error:	エラー
func (h checkInOutHandler) GetTodayCheckins(c echo.Context) ([]dto.CheckinsRes, error) {
	// ログインユーザーIDの取得
	dogownerID, err := wrcontext.GetLoginDogownerID(c)
	if err != nil {
		return nil, err
	}
	//dogownerの今日のチェックイン履歴の取得
	checkinsResult, err := h.r.GetTodayCheckinsByDogownerID(c, dogownerID)
	if err != nil {
		return nil, err
	}

	//検索結果をレスポンスに詰める
	checkinsRes := []dto.CheckinsRes{}
	for _, checkinResult := range checkinsResult {
		checkinRes := dto.CheckinsRes{
			DogID:       checkinResult.DogID.Int64,
			DogrunID:    checkinResult.DogrunID.Int64,
			CheckinAt:   checkinResult.CheckinAt.Time,
			ReCheckinAt: checkinResult.ReCheckinAt.Time,
		}
		checkinsRes = append(checkinsRes, checkinRes)
	}
	return checkinsRes, nil
}
