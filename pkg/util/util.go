package util

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"strings"
	"time"
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

// UUIDGenerator: 指定された長さの一意のIDを生成
// Args:
//
//	length: 生成されるIDの長さを指定
//	handleError: エラーを受け取りエラーを返す、カスタムエラーを処理する関数
//
// Returns:
//
//	string: ユニークなID
//	error: error情報
func UUIDGenerator(l int, handleError func(error) error) (string, error) {
	b := make([]byte, l)
	_, err := rand.Read(b)

	if err != nil {
		return "", handleError(err)
	}

	return base64.RawURLEncoding.EncodeToString(b)[:l], nil
}
