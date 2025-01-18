package common

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

// PKカスタムバリデーション（登録時）
func VCreatePrimaryKey(fl validator.FieldLevel) bool {
	pk := fl.Field().Int()
	return pk == 0
}

// PKカスタムバリデーション（更新時）
func VUpdatePrimaryKey(fl validator.FieldLevel) bool {
	pk := fl.Field().Int()
	return pk != 0
}

// 性別値
const (
	SEX_MALE    = "M"
	SEX_FEMALE  = "F"
	SEX_UNKNOWN = "U"
	SEX_OTHER   = "O"
)

var sex_values []string = []string{SEX_MALE, SEX_FEMALE, SEX_UNKNOWN, SEX_OTHER}

// 性別のバリデーションs
func VSex(fl validator.FieldLevel) bool {
	sex := fl.Field().String()

	for _, v := range sex_values {
		if v == sex {
			return true
		}
	}
	return false
}

// スライスの空バリデーション
func VNotEmpty(fl validator.FieldLevel) bool {
	slice := fl.Field()

	field := fl.Field()

	// フィールドがスライスかどうかを確認
	if field.Kind() == reflect.Slice {
		return slice.Len() > 0 // スライスの長さが0より大きければOK
	}
	return false // スライス以外は無効
}

// int64スライスの最大値
