package dto

// bookmark 登録用
type BookmarkAddReq struct {
	DogrunIDs []int64 `json:"bookmark_dogrun_id" validate:"required,notEmpty"`
}

type BookmarkDeleteReq struct {
	DogrunIDs []int64 `json:"bookmark_dogrun_id" validate:"required,notEmpty"`
}

type CheckinReq struct {
	DogrunID int64   `json:"dogrun_id" validate:"required"`
	DogIDs   []int64 `json:"dog_id" validate:"required"`
}

type CheckoutReq struct {
	DogrunID int64   `json:"dogrun_id" validate:"required"`
	DogIDs   []int64 `json:"dog_id" validate:"required"`
}
