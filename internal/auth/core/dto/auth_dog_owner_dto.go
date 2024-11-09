package dto

type AuthDogOwnerReq struct {
	Password          string `json:"password"`
	DogOwnerName      string `json:"dogOwnerName"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phoneNumber"`
	AuthorizationCode string
}
