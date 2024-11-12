package dto

import (
	"github.com/wanrun-develop/wanrun/common"
)

// dog詳細レスポンス
type DogDetailsRes struct {
	DogID      int64         `json:"dogId"`
	DogOwnerID int64         `json:"dogOwnerId"`
	Name       string        `json:"name"`
	Weight     int64         `json:"weight"`
	Sex        string        `json:"sex"`
	Image      string        `json:"image"`
	DogTypeId  []int64       `json:"dogTypeId"`
	CreateAt   common.WRTime `json:"createAt"`
	UpdateAt   common.WRTime `json:"updateAt"`
}

// dog一覧用レスポンス
type DogListRes struct {
	DogID     int64   `json:"dogId"`
	Name      string  `json:"name"`
	Weight    int64   `json:"weight"`
	Sex       string  `json:"sex"`
	Image     string  `json:"image"`
	DogTypeId []int64 `json:"dogTypeId"`
}

// dogType用レスポンス
type DogTypeMstRes struct {
	DogTypeID int    `json:"dogTypeId"`
	Name      string `json:"name"`
}
