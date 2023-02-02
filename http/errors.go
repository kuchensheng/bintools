package http

import "encoding/json"

type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

var (
	UNKNOWN   = BusinessError{20000, "未知异常", nil}
	SERVERERR = BusinessError{20500, "服务端异常", nil}
)

func NewError(data any) BusinessError {
	return BusinessError{-1, "服务端错误", data}
}

func NewErrorWithCode(code int, data any) BusinessError {
	return BusinessError{code, "服务端错误", data}
}

func NewErrorWithMsg(code int, message string, data any) BusinessError {
	return BusinessError{code, message, data}
}

func (e *BusinessError) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
