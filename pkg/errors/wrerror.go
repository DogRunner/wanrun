package errors

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type wrError struct {
	eType      eType
	causeBy    string
	msg        string
	innerError error
}

func (me *wrError) Error() string {
	return fmt.Sprintf("wanrun error: code[%s], message[%s]", me.eType, me.msg)
}

func (e wrError) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		fmt.Fprintf(f, "  code[%s]", e.eType)
		fmt.Fprintf(f, "  message: %s", e.msg)
		fmt.Fprintf(f, " ->causeBy: %s", e.causeBy)
	default:
		fmt.Fprintf(f, "%v", string(c))
	}
}

/*
エラー生成
すでにWRErrorの場合は、そのまま返す
*/
func NewWRError(err error, msg string, errorType eType) *wrError {
	// errがnilの場合は根本エラーは指定しない
	if err == nil {
		return &wrError{
			eType: errorType,
			msg:   msg,
		}
	}
	//wrError型にアサーションできるかチェック
	if me, ok := err.(*wrError); ok {
		return me // 既存のwrErrorを返す
	}

	return &wrError{
		eType:      errorType,
		causeBy:    err.Error(),
		msg:        msg,
		innerError: err,
	}
}

/*
カスタムエラーハンドラーミドルウェア
*/
func HttpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError

	var me *wrError
	if wreer, ok := err.(*wrError); ok {
		me = wreer
		code = mappingError(me)
	} else {
		c.Logger().Error(err)
		_ = c.JSON(code, err)
		return
	}

	res := ErrorRes{
		Code:       me.eType.String(),
		Message:    me.msg,
		StackTrace: me.causeBy,
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				c.Logger().Error(err)
			}
		} else {
			err := c.JSON(code, res)
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}

/*
内部エラーコードどHTTPエラーコードのマッピング
*/
func mappingError(err *wrError) int {
	errorType := err.eType

	var httpCode int
	switch errorType.errorType {
	case CLIENT:
		switch errorType.service {
		case AUTH:
			httpCode = http.StatusUnauthorized //401
		default:
			httpCode = http.StatusBadRequest //400
		}
	case SERVER:
		httpCode = http.StatusInternalServerError //500
	default:
		httpCode = http.StatusInternalServerError //500
	}

	return httpCode
}
