package errors

type ErrorResponse struct {
	Message    string `json:"message"`
	StackTrace string `json:"trace,omitempty"`
}

type ErrorRes struct {
	Message    string `json:"message"`
	StackTrace string `json:"trace,omitempty"`
}
