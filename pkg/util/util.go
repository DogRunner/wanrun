package util

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

func IsStrEmpty(s string) bool {
	// strings.TrimSpaceで空白を取り除き、空文字をチェック
	return strings.TrimSpace(s) == ""
}

func IsPstrEmpty(s *string) bool {
	if s == nil {
		return true
	}
	// strings.TrimSpaceで空白を取り除き、空文字をチェック
	return strings.TrimSpace(*s) == ""
}

/*
HH:mm:ssをtime.Timeに変換。エラーの場合は初期値を返す
*/
func ParseStrToTime(str string) time.Time {
	t, err := time.Parse("15:04:05", str)
	if err != nil {
		return time.Time{}
	}
	return t
}

/*
sql.NullStringのバリデーション。
valid = falseの際は、デフォルト値を返す
*/
func ChooseStringValidValue(sqlStr sql.NullString, str string) string {
	if sqlStr.Valid {
		return sqlStr.String
	}
	return str
}

/*
sql.NullInt64のバリデーション。
valid = falseの際は、デフォルト値を返す
*/
func ChooseInt64ValidValue(sqlInt sql.NullInt64, i int64) int64 {
	if sqlInt.Valid {
		return sqlInt.Int64
	}
	return i
}

/*
sql.NullFloat64のバリデーション。
valid = falseの際は、デフォルト値を返す
*/
func ChooseFloat64ValidValue(sqlFloat sql.NullFloat64, f float64) float64 {
	if sqlFloat.Valid {
		return sqlFloat.Float64
	}
	return f
}

/*
sql.NullBoolのバリデーション。
valid = falseの際は、デフォルト値を返す
*/
func ChooseBoolValidValue(sqlBool sql.NullBool, b bool) bool {
	if sqlBool.Valid {
		return sqlBool.Bool
	}
	return b
}

/*
t time.のバリデーション。
valid = falseの際は、デフォルト値を返す
*/
func ChooseTimeValidValue(sqlTime sql.NullTime, t time.Time) time.Time {
	if sqlTime.Valid {
		return sqlTime.Time
	}
	return t
}

// ConvertSliceToMap: sliceをmapに変換するutil関数
//
// args:
//   - []T:	mapにしたいスライス
//   - func(T) K:	mapのkeyを抽出する関数
//
// return:
//   - map[K]T:	変換したmap
func ConvertSliceToMap[T any, K comparable](slice []T, keySelector func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range slice {
		key := keySelector(item)
		result[key] = item
	}
	return result
}

// UUIDGenerator: UUIDを生成する
// Args:
//
//	handleError: エラーを処理するカスタム関数
//
// Returns:
//
//	string: UUID
//	error: エラー情報
func UUIDGenerator(handleError func(error) error) (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", handleError(err)
	}
	return u.String(), nil
}

// ConvertStringPointer: awsなどで返ってくる*string型をstringに返す
// Args:
//
//	*string: stringのポインター型
//
// Returns:
//
//	string: string型の文字列
func ConvertStringPointer(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
