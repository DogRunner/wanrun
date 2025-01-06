package dto

type AuthDogownerReq struct {
	Password          string `json:"password"`
	DogownerName      string `json:"dogOwnerName"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phoneNumber"`
	AuthorizationCode string
}
