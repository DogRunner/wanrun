package dto

type DogOwnerReq struct {
	Password     string `json:"password"`
	DogOwnerName string `json:"dogOwnerName"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
}

type DogOwnerDTO struct {
	DogOwnerID int64
	JwtID      string
}
