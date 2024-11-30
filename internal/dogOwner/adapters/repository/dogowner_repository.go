package repository

import (
	"github.com/labstack/echo/v4"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type IDogOwnerRepository interface {
	GetDogOwnerById(int64) (model.DogOwner, error)
	CreateDogOwner(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error
}

type dogOwnerRepository struct {
	db *gorm.DB
}

func NewDogRepository(db *gorm.DB) IDogOwnerRepository {
	return &dogOwnerRepository{db}
}

func (dr *dogOwnerRepository) GetDogOwnerById(dogOwnerId int64) (model.DogOwner, error) {
	dogOwner := model.DogOwner{}
	if err := dr.db.Where("dog_owner_id = ?", dogOwnerId).Find(&dogOwner).Error; err != nil {
		return model.DogOwner{}, err
	}
	return dogOwner, nil
}

// CreateDogOwner: DogOwnerの作成
//
// args:
//   - echo.Context: c Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - *model.DogOwnerCredential: doc ドッグオーナーのクレデンシャル
//
// return:
//   - *model.DogOwnerCredential: ドッグオーナーのクレデンシャル
//   - error: error情報
func (dor *dogOwnerRepository) CreateDogOwner(tx *gorm.DB, c echo.Context, doc *model.DogOwnerCredential) error {
	logger := log.GetLogger(c).Sugar()

	// dog_ownersテーブルにDogOwnerの作成
	if err := tx.Create(&doc.AuthDogOwner.DogOwner).Error; err != nil {
		logger.Error("Failed to create DogOwner: ", err)
		return wrErrors.NewWRError(
			err,
			"DogOwner作成に失敗しました。",
			wrErrors.NewDogOwnerServerErrorEType(),
		)
	}

	// DogOwnerが作成された後、そのIDをauthDogOwnerに設定
	doc.AuthDogOwner.DogOwnerID = doc.AuthDogOwner.DogOwner.DogOwnerID

	logger.Infof("Created DogOwner Detail: %v", doc.AuthDogOwner.DogOwner)

	return nil
}
