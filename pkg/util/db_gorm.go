package util

import (
	"database/sql"
	"time"

	"github.com/wanrun-develop/wanrun/common"
)

/*
時間型用の構造体
*/
type CustomTime struct {
	sql.NullTime
}

/*
時間のみの構造体に変換
gormはpsqlのtime型を自動でキャストできないようなので、実装
*/
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Valid = false
		return nil
	}
	ct.Valid = true
	var t time.Time
	switch v := value.(type) {
	case string:
		var err error
		t, err = time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
	case time.Time:
		t = v
	default:
		ct.Valid = false
		return nil
	}
	ct.Time = t
	return nil
}

/*
string型の値をsql.NullStringに変換
*/
func NewSqlNullString(value string) sql.NullString {
	// 値が空文字またはnullの場合は、`Valid: false` と設定し、`NULL` として扱う
	if value == "" || value == "null" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	// 値が存在する場合は `Valid: true` として設定
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

/*
int64型の値をsql.NullInt64に変換
*/
func NewSqlNullInt64(value int64) sql.NullInt64 {
	// 値がゼロまたは負の値など、無効な条件があれば `Valid: false` と設定
	if value == 0 {
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}
	// 有効な値の場合は `Valid: true` として設定
	return sql.NullInt64{
		Int64: value,
		Valid: true,
	}
}

/*
bool型の値をsql.NullBoolに変換
*/
func NewSqlNullBool(value bool) sql.NullBool {
	return sql.NullBool{
		Bool:  value,
		Valid: true,
	}
}

/*
time.Time型の値をsql.NullTimeに変換
*/
func NewSqlNullTime(value time.Time) sql.NullTime {
	// `value` がゼロ値の場合は `Valid: false` と設定
	if value.IsZero() {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}
	// `value` が有効な時刻であれば `Valid: true` として設定
	return sql.NullTime{
		Time:  value,
		Valid: true,
	}
}

/*
time.Time型の値をCustomTimeに変換
*/
func NewCustomTime(value time.Time) CustomTime {
	// `value` がゼロ値かどうかを確認し、`Valid` フィールドを設定
	if value.IsZero() {
		return CustomTime{
			NullTime: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		}
	}
	return CustomTime{
		NullTime: sql.NullTime{
			Time:  value,
			Valid: true,
		},
	}
}

/*
sql.NullTime型の値をCustomTimeに変換
*/
func NewCustomTimeFromNullTime(nullTime sql.NullTime) CustomTime {
	return CustomTime{NullTime: nullTime}
}

// ConvetToWRTime: sql.NullTimeをカスタムTime型に変換
//
// args:
//   - sql.NullTime:	変換対象
//
// return:
//   - common.WRTime:	カスタム時間
//   - :
func ConvertToWRTime(nulltime sql.NullTime) common.WRTime {
	return common.WRTime{Time: nulltime.Time}
}

/*
byte型の値をsql.NullByteに変換
*/
func NewSqlNullByte(value []byte) sql.NullByte {
	if len(value) == 0 {
		return sql.NullByte{
			Byte:  0, // デフォルトの0バイトを使用
			Valid: false,
		}
	}
	return sql.NullByte{
		Byte:  value[0], // TODO: 最初のバイトを取得だけど他に良い方法が欲しい。
		Valid: true,
	}
}
