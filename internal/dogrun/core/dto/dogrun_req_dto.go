package dto

/*
円型検索のリクエストボディ
*/
type SearchAroundCircleCondition struct {
	Center pointer      `json:"center" validate:"required"`
	Target centerTarget `json:"target" validate:"required"`
}

/*
検索対象の中心点
*/
type pointer struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`   //経度
	Longitude float64 `json:"longitude" validate:"required,longitude"` //緯度
}

/*
検索対象の距離
単位: m(メートル)
*/
type centerTarget struct {
	Radius int `json:"radius" validate:"required,gte=0,lte=50000"` // 半径（0以上）
}

/*
長方形型検索のリクエストボディ
*/
type SearchAroundRectangleCondition struct {
	Target            rectangleTarget `json:"target" validate:"required"`
	IncludeDogrunTags []int64         `json:"include_dogrunTags" validate:"min=0,max=100"`
}

/*
長方形の南西（右下）と北東（右上）を示す
*/
type rectangleTarget struct {
	Southwest pointer `json:"southwest" validate:"required"`
	Northeast pointer `json:"northeast" validate:"required"`
}
