package dto

type OrgReq struct {
	OrgName      string `json:"organizationName" validate:"required"`
	ContactEmail string `json:"contactEmail" validate:"required"`
	PhoneNumber  string `json:"phoneNumber" validate:"required"`
	Address      string `json:"address" validate:"required"`
	Description  string `json:"description"`
	Password     string `json:"password" validate:"required"`
}
