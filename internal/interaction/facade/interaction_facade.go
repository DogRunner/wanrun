package facade

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/interaction/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/wrcontext"
)

type IBookmarkFacade interface {
	GetAllUserBookmarks(echo.Context) ([]int64, error)
}

type bookmarkFacade struct {
	r repository.IBookmarkRepository
}

func NewBookmarkFacade(br repository.IBookmarkRepository) IBookmarkFacade {
	return &bookmarkFacade{br}
}

// GetAllUserBookmarks: ログインユーザーのブックマークを取得
//
// args:
//   - echo.Context:	コンテキスト
//
// return:
//   - []int64:	bookmarkIDs
//   - error:	エラー
func (f *bookmarkFacade) GetAllUserBookmarks(c echo.Context) ([]int64, error) {
	// ログインユーザーIDの取得
	userID, err := wrcontext.GetLoginUserID(c)
	if err != nil {
		return nil, err
	}

	bookmarks, err := f.r.GetBookmarks(c, userID)
	if err != nil {
		return nil, err
	}

	bookmarkedDogrunIDs := []int64{}
	for _, bookmark := range bookmarks {
		bookmarkedDogrunIDs = append(bookmarkedDogrunIDs, bookmark.DogrunID.Int64)
	}

	return bookmarkedDogrunIDs, nil
}
