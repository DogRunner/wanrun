package dto

// bookmark 登録用
type AddBookmark struct {
	DogrunIDs []int64 `json:"bookmarkDogrunId" validate:"required,notEmpty"`
}

type DeleteBookmark struct {
	DogrunIDs []int64 `json:"bookmarkDogrunId" validate:"required,notEmpty"`
}
