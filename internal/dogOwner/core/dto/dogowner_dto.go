package dto

type DogOwnerReq struct {
	Password     string `json:"password" validate:"required"`
	DogOwnerName string `json:"dogOwnerName" validate:"required"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
}

type DogOwnerDTO struct {
	DogOwnerID int64
	JwtID      string
}
