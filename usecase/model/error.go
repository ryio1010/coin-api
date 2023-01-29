package model

type ErrorResponse struct {
	ErrorCode int    `json:"errorcode"`
	Message   string `json:"message"`
}

func CreateResponse(code int, msg string) ErrorResponse {
	e := ErrorResponse{
		ErrorCode: code,
		Message:   msg,
	}
	return e
}
