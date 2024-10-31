package dto

import "time"

type DogDetailsRes struct {
	DogID      int64     `json:"dogId"`
	DogOwnerID int64     `json:"dogOwnerId"`
	Name       string    `json:"name"`
	DogTypeID  int64     `json:"dogTypeId"`
	Weight     int64     `json:"weight"`
	Sex        string    `json:"sex"`
	Image      string    `json:"image"`
	CreateAt   time.Time `json:"createAt"`
	UpdateAt   time.Time `json:"updateAt"`
}

type DogListRes struct {
	DogID     int64  `json:"dogId"`
	Name      string `json:"name"`
	DogTypeID int64  `json:"dogTypeId"`
	Weight    int64  `json:"weight"`
	Sex       string `json:"sex"`
	Image     string `json:"image"`
}

type DogTypeMstRes struct {
	DogTypeID int    `json:"dogTypeId"`
	Name      string `json:"name"`
}
