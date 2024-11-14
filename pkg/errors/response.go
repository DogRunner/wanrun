package errors

type ErrorResponse struct {
	Message    string `json:"message"`
	StackTrace string `json:"trace,omitempty"`
}

type ErrorRes struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StackTrace string `json:"trace,omitempty"`
}

// / NewErrorRes: ErrorResの生成
//
// args:
//   - *wrError: wanrunエラー
//
// return:
//   - ErrorRes: エラーレスポンス情報
func NewErrorRes(we error) ErrorRes {
	me, ok := we.(*wrError)
	if !ok {
		me = NewWRError(
			we,
			"予期せぬエラーが起きました。",
			NewUnexpectedErrorEType(),
		)
	}
	return ErrorRes{
		Code:       me.eType.String(),
		Message:    me.msg,
		StackTrace: me.causeBy,
	}
}
