package dto

type DogownerReq struct {
	Password     string `json:"password" validate:"required"`
	DogownerName string `json:"dogOwnerName" validate:"required"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
}
