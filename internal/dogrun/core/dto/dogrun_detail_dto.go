package dto

import "time"

// ドッグラン詳細情報用
type DogrunDetailDto struct {
	DogrunID        int            `json:"dogrun_id,omitempty"`
	DogrunManagerID int            `json:"dogrun_manager_id,omitempty"`
	PlaceId         string         `json:"place_id,omitempty"`
	Name            string         `json:"name"`
	Address         Address        `json:"address"`
	Location        Location       `json:"location"`
	BusinessStatus  string         `json:"business_status"`
	NowOpen         bool           `json:"now_open"`
	BusinessDay     int            `json:"business_day,omitempty"`
	Holiday         int            `json:"holiday,omitempty"`
	OpenTime        string         `json:"open_time"`
	CloseTime       string         `json:"close_time"`
	Description     string         `json:"description,omitempty"`
	GoogleRating    float32        `json:"google_rating,omitempty"`
	UserRatingCount int            `json:"user_rating_count,omitempty"`
	DogrunTags      []DogrunTagDto `json:"dogrun_tags,omitempty"`
	CreateAt        *time.Time     `json:"create_at,omitempty"`
	UpdateAt        *time.Time     `json:"update_at,omitempty"`
}

// ドッグランタグ情報
type DogrunTagDto struct {
	DogrunTagId int    `json:"dogrun_tag_id"`
	TagId       int    `json:"tag_id"`
	TagName     string `json:"tag_idag_name"`
	Description string `json:"description"`
}

// 軽度・緯度情報
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// 住所情報
type Address struct {
	PostCode string `json:"postcode"`
	Address  string `json:"address"`
}
