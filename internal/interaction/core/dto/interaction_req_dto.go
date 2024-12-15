package dto

// bookmark 登録用
type BookmarkAddReq struct {
	DogrunIDs []int64 `json:"bookmarkDogrunId" validate:"required,notEmpty"`
}

type BookmarkDeleteReq struct {
	DogrunIDs []int64 `json:"bookmarkDogrunId" validate:"required,notEmpty"`
}
