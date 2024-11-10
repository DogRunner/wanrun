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
