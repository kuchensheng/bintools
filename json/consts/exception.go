package consts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

var ErrKey = "_err_key"

type exception struct {
	Location    string `json:"location"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (e *exception) Error() string {
	return e.Description
}
func NewException(location, name, desc string) *exception {
	return &exception{location, name, desc}
}

func DeferHandler() error {
	if x := recover(); x != nil {
		log.Error().Msgf("发生了panic错误:%v", x.(error))
		return x.(error)
	}
	return nil
}

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

func Ok() *BusinessException {
	return &BusinessException{
		Code:    0,
		Message: "成功",
	}
}

func OkWithData(data any) *BusinessException {
	e := Ok()
	e.Data = data
	return e
}
