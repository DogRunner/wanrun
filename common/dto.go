package common

import (
	"fmt"
	"time"
)

type WRTime struct {
	time.Time
}

// フォーマットを yyyy/MM/dd HH:mm:ss に固定
const F_yyyyMMddHHmmss = "2006/01/02 15:04:05"

// WRTime の JSON 出力用メソッドをカスタマイズ
func (ct WRTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", ct.Format(F_yyyyMMddHHmmss))
	return []byte(formatted), nil
}
