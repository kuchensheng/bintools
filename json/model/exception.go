package model

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
)

type BusinessException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (exception *BusinessException) Error() string {
	return exception.Message
}

func NewBusinessException(code int, message string) *BusinessException {
	return &BusinessException{
		Code:    code,
		Message: message,
	}
}

func NewBusinessExceptionWithData(code int, message string, data any) *BusinessException {
	return &BusinessException{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (exception *BusinessException) Write2Response(request *http.Request, statusCode int) {
	request.Response.StatusCode = statusCode
	if data, err := json.Marshal(exception); err == nil {
		request.Response.Write(bufio.NewWriter(bytes.NewBuffer(data)))
	} else {
		request.Response.Write(bufio.NewWriter(bytes.NewBufferString(err.Error())))
	}
}

func (exception *BusinessException) WriteResponse(writer http.ResponseWriter, request *http.Request, statusCode int) {
	request.Response.StatusCode = statusCode
	if data, err := json.Marshal(exception); err == nil {
		writer.Write(data)
	} else {
		writer.Write([]byte(err.Error()))
	}
}
