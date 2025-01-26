package dto

type DogrunmgReq struct {
	Password     string `json:"password" validate:"required"`
	DogrunmgName string `json:"dogrunmgName" validate:"required"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
}
