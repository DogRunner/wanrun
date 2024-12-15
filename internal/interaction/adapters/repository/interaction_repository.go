package repository

import (
	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
	"gorm.io/gorm"
)

type IBookmarkRepository interface {
	GetBookmarks(echo.Context, int64) ([]model.DogrunBookmark, error)
	AddBookmark(echo.Context, int64, int64) (int64, error)
	FindDogrunBookmark(echo.Context, int64, int64) (model.DogrunBookmark, error)
	DeleteBookmark(echo.Context, []int64, int64) error
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) IBookmarkRepository {
	return &bookmarkRepository{db}
}

// GetBookmarks: dogownerのブックマークを取得
//
// args:
//   - echo.Context:	コンテキスト
//   - int64:	ドッグオーナーID
//
// return:
//   - []model.DogrunBookmark:	検索結果
//   - error:	エラー
func (r *bookmarkRepository) GetBookmarks(c echo.Context, dogownerID int64) ([]model.DogrunBookmark, error) {
	logger := log.GetLogger(c).Sugar()

	bookmarks := []model.DogrunBookmark{}
	if err := r.db.
		Where("dog_owner_id = ?", dogownerID).
		Find(&bookmarks).Error; err != nil {
		logger.Error(err)
		err := errors.NewWRError(err, "dogrun_bookmarksの検索に失敗しました。", errors.NewInteractionServerErrorEType())
		return nil, err
	}

	return bookmarks, nil
}

// AddBookmark: ドックランのブックマーク登録
//
// args:
//   - echo.Context:	コンテキスト
//   - int64:	ブックマーク対象のdogrunId
//   - int64:	ブックマーク登録者のdogownerId
//
// return:
//   - int64:	発行されたdogrun_bookmark_id
//   - error:	エラー
func (r *bookmarkRepository) AddBookmark(c echo.Context, dogrunID int64, dogownerID int64) (int64, error) {
	logger := log.GetLogger(c).Sugar()
	bookmark := model.DogrunBookmark{
		DogOwnerID: util.NewSqlNullInt64(dogownerID),
		DogrunID:   util.NewSqlNullInt64(dogrunID),
	}

	if err := r.db.Create(&bookmark).Error; err != nil {
		logger.Error(err)
		err := errors.NewWRError(err, "dogrun_bookmarkの登録に失敗しました。", errors.NewInteractionServerErrorEType())
		return 0, err
	}
	return bookmark.DogrunBookmarkID.Int64, nil
}

// FindDogrunBookmark: dogrunIdとdogownerIdでbookmarkへ検索
//
// args:
//   - echo.Context:	コンテキスト
//   - int64:	dogrunIDで条件指定
//   - int64:	dogownerIDで条件指定
//
// return:
//   - model.DogrunBookmark:	検索結果構造体
//   - error:	エラー
func (r *bookmarkRepository) FindDogrunBookmark(c echo.Context, dogrunID int64, dogownerID int64) (model.DogrunBookmark, error) {
	logger := log.GetLogger(c).Sugar()

	bookmark := model.DogrunBookmark{}
	if err := r.db.
		Where("dogrun_id = ?", dogrunID).
		Where("dog_owner_id = ?", dogownerID).
		Find(&bookmark).Error; err != nil {
		logger.Error(err)
		err := errors.NewWRError(err, "dogrun_bookmarkの検索に失敗しました。", errors.NewInteractionServerErrorEType())
		return bookmark, err
	}

	return bookmark, nil
}

// AddBookmark: 複数ドックランのブックマーク削除
//
// args:
//   - echo.Context:	コンテキスト
//   - []int64:	ブックマーク削除対象のdogrunIds
//   - int64:	ブックマーク登録者のdogownerId
//
// return:
//   - error:	エラー
func (r *bookmarkRepository) DeleteBookmark(c echo.Context, dogrunIDs []int64, dogownerID int64) error {

	logger := log.GetLogger(c).Sugar()

	if err := r.db.
		Where("dog_owner_id = ?", dogownerID).
		Where("dogrun_id IN ?", dogrunIDs).
		Delete(&model.DogrunBookmark{}).Error; err != nil {
		logger.Error(err)
		err := errors.NewWRError(err, "dogrun_bookmarkの削除に失敗しました。", errors.NewInteractionServerErrorEType())
		return err
	}
	return nil
}