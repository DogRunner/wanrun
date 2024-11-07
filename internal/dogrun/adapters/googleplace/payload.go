package googleplace

import (
	"github.com/wanrun-develop/wanrun/internal/dogrun/core/dto"
)

const (
	RANKPREFERENCE_POPULARITY = "POPULARITY" //人気度に基づいて検索結果を並べ替えます。
	RANKPREFERENCE_DISTANCE   = "DISTANCE"   //結果の距離に基づいて昇順で並べ替えます
	RANKPREFERENCE_RELEVANCE  = "RELEVANCE"  //テキストクエリ結果を検索関連性に基づいて並べ替えます。
)

// https://developers.google.com/maps/documentation/places/web-service/nearby-search
type SearchNearbyPayLoad struct {
	IncludedTypes       []string                  `json:"includedTypes" validate:"required"`               // 含む場所タイプ　https://developers.google.com/maps/documentation/places/web-service/place-types?hl=ja#table-a
	ExcludedTypes       []string                  `json:"excludedTypes,omitempty"`                         //除く場所タイプ
	MaxResultCount      int                       `json:"maxResultCount" validate:"required,gte=1,lte=20"` // 最大結果数（1以上）
	LocationRestriction locationRestrictionCircle `json:"locationRestriction" validate:"required"`         // 位置制限
	RankPreference      string                    `json:"rankPreference,omitempty"`                        //検索のランキング付け
}

// https://developers.google.com/maps/documentation/places/web-service/text-search
type SearchTextPayLoad struct {
	TextQuery           string                       `json:"textQuery" validate:"required"`             // 検索文言
	IncludedTypes       []string                     `json:"includedTypes,omitempty"`                   // 含む場所タイプ　https://developers.google.com/maps/documentation/places/web-service/place-types?hl=ja#table-a
	ExcludedTypes       []string                     `json:"excludedTypes,omitempty"`                   // 除く場所タイプ
	PageSize            int                          `json:"pageSize" validate:"required,gte=1,lte=20"` // 最大結果数（1以上）
	LocationRestriction locationRestrictionRectangle `json:"locationRestriction" validate:"required"`   // 位置制限
	RankPreference      string                       `json:"rankPreference,omitempty"`                  // 検索のランキング付け
	PageToken           string                       `json:"pageToken,omitempty"`                       //次回ページのトークン
}

type locationRestrictionCircle struct {
	Circle circle `json:"circle" validate:"required"` // サークル型の位置制限
}

type circle struct {
	Center pointer `json:"center" validate:"required"`                 // 中心点
	Radius float64 `json:"radius" validate:"required,gte=0,lte=50000"` // 半径（0以上）
}

type pointer struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`   // 緯度
	Longitude float64 `json:"longitude" validate:"required,longitude"` // 経度
}

type locationRestrictionRectangle struct {
	Rectangle rectangle `json:"rectangle" validate:"required"` //長方形型の一条号
}

type rectangle struct {
	Low  pointer `json:"low" validate:"required"`  //長方形の南西の角
	High pointer `json:"high" validate:"required"` //長方形の北東の角
}

func ConvertReqToSearchNearbyPayload(req dto.SearchAroundCircleCondition) SearchNearbyPayLoad {
	center := pointer{
		Latitude:  req.Center.Latitude,
		Longitude: req.Center.Longitude,
	}

	circle := circle{
		Center: center,
		Radius: float64(req.Target.Radius),
	}

	return SearchNearbyPayLoad{
		IncludedTypes:       []string{"dog_park"},
		MaxResultCount:      10,
		LocationRestriction: locationRestrictionCircle{Circle: circle},
		RankPreference:      "POPULARITY",
	}
}

func ConvertReqToSearchTextPayload(req dto.SearchAroundRectangleCondition) SearchTextPayLoad {
	rectangle := rectangle{
		Low: pointer{
			Latitude:  req.Target.Southwest.Latitude,
			Longitude: req.Target.Southwest.Longitude,
		},
		High: pointer{
			Latitude:  req.Target.Northeast.Latitude,
			Longitude: req.Target.Northeast.Longitude,
		},
	}
	locationRestrictionRectangle := locationRestrictionRectangle{
		Rectangle: rectangle,
	}

	return SearchTextPayLoad{
		TextQuery:           "ドッグラン",
		PageSize:            20,
		LocationRestriction: locationRestrictionRectangle,
		RankPreference:      RANKPREFERENCE_RELEVANCE,
	}
}
