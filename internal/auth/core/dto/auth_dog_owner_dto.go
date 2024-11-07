package dto

type ReqAuthDogOwnerDto struct {
	Password          string `json:"password"`
	DogOwnerName      string `json:"dogOwnerName"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phoneNumber"`
	AuthorizationCode string
}

type ResAuthDogOwner struct {
	AuthDogOwnerID int64  `json:"authDogOwnerId"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Sex            string `json:"sex"`
}
