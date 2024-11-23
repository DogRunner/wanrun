package dto

// bookmark 登録用
type AddBookmark struct {
	DogrunID int64 `json:"bookmarkDogrunId" validate:"required"`
}
