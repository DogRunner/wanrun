package dto

// dogのsave用
type DogSaveReq struct {
	DogID      int64  `json:"dogId" validate:"primaryKey"`
	DogownerID int64  `json:"dogOwnerId" validate:"required"`
	Name       string `json:"name" validate:"required"`
	DogTypeID  int64  `json:"dogTypeId" validate:"required"`
	Weight     int64  `json:"weight" validate:"required"`
	Sex        string `json:"sex" validate:"required,sex"`
	Image      string `json:"image"`
}
