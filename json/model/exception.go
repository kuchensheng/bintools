package model

import "encoding/json"

type BusinessException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (exception *BusinessException) Error() string {
	if marshal, err := json.Marshal(exception); err != nil {
		return err.Error()
	} else {
		return string(marshal)
	}
}

func NewBusinessException(code int, message string) *BusinessException {
	return &BusinessException{
		Code:    code,
		Message: message,
	}
}
