package dto

import "time"

// ドッグラン詳細画面での表示情報
type DogrunDetail struct {
	DogrunID        int            `json:"dogrunId,omitempty"`
	DogrunManagerID int            `json:"dogrunManagerId,omitempty"`
	PlaceId         string         `json:"placeId,omitempty"`
	Name            string         `json:"name"`
	Address         Address        `json:"address"`
	Location        Location       `json:"location"`
	BusinessStatus  string         `json:"businessStatus,omitempty"`
	NowOpen         bool           `json:"nowOpen"`
	BusinessHour    BusinessHour   `json:"businessHour"`
	Description     string         `json:"description,omitempty"`
	GoogleRating    float32        `json:"googleRating,omitempty"`
	UserRatingCount int            `json:"userRatingCount,omitempty"`
	DogrunTags      []DogrunTagDto `json:"dogrunTags,omitempty"`
	CreateAt        *time.Time     `json:"createAt,omitempty"`
	UpdateAt        *time.Time     `json:"updateAt,omitempty"`
}

// ドッグラン一覧での表示情報
type DogrunLists struct {
	DogrunID          int             `json:"dogrunId,omitempty"`
	PlaceId           string          `json:"placeId,omitempty"`
	Name              string          `json:"name"`
	Address           Address         `json:"address"`
	Location          Location        `json:"location"`
	BusinessStatus    string          `json:"businessStatus,omitempty"`
	NowOpen           bool            `json:"nowOpen"`
	ToadyBusinessHour DayBusinessTime `json:"toadyBusinessHour"`
	Description       string          `json:"description,omitempty"`
	GoogleRating      float32         `json:"googleRating,omitempty"`
	UserRatingCount   int             `json:"userRatingCount,omitempty"`
	DogrunTags        []DogrunTagDto  `json:"dogrunTags,omitempty"`
	Photos            []PhotoInfo     `json:"photos,omitempty"`
}

// 営業日情報
type BusinessHour struct {
	Regular RegularBusinessHour   `json:"regular"`
	Special []SpecialBusinessHour `json:"special,omitempty"`
}

// 通常営業日情報
type RegularBusinessHour struct {
	Sunday    DayBusinessTime `json:"0,omitempty"`
	Monday    DayBusinessTime `json:"1,omitempty"`
	Tuesday   DayBusinessTime `json:"2,omitempty"`
	Wednesday DayBusinessTime `json:"3,omitempty"`
	Thursday  DayBusinessTime `json:"4,omitempty"`
	Friday    DayBusinessTime `json:"5,omitempty"`
	Saturday  DayBusinessTime `json:"6,omitempty"`
}

type DayBusinessTime struct {
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
	IsAllDay  bool   `json:"isAllDay"`
	IsHoliday bool   `json:"isHoliday"`
}

// 特別営業日情報
type SpecialBusinessHour struct {
	Date string `json:"date"`
	DayBusinessTime
}

type PhotoInfo struct {
	PhotoKey string `json:"photoKey"`
	WidthPx  uint   `json:"widthPx"`
	HeightPx uint   `json:"heightPx"`
}

// ドッグランタグ情報
type DogrunTagDto struct {
	DogrunTagId int    `json:"dogrunTagId"`
	TagId       int    `json:"tagId"`
	TagName     string `json:"tagIdagName"`
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
