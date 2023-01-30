package model

type ErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func CreateErrorResponse(code int, msg string) *ErrorResponse {
	e := &ErrorResponse{
		ErrorCode: code,
		Message:   msg,
	}
	return e
}
