package http

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

type HttpResponse struct {
	response *http.Response
}

func (resp *HttpResponse) GetBytes() []byte {
	response := resp.response
	if data, err := io.ReadAll(response.Body); err != nil {
		log.Error("无法读取响应体内容:%v", err)
		return nil
	} else {
		return data
	}
}

func (resp *HttpResponse) GetString() string {
	return string(resp.GetBytes())
}

func (resp *HttpResponse) GetContentLength() int64 {
	return resp.response.ContentLength
}

func (resp *HttpResponse) GetStatus() int {
	return resp.response.StatusCode
}

func (resp *HttpResponse) IsOk() bool {
	return resp.GetStatus() >= http.StatusOK && resp.GetStatus() <= http.StatusMultipleChoices
}

func (resp *HttpResponse) WriteString(content string) error {
	return resp.response.Write(bytes.NewBufferString(content))
}

func (resp *HttpResponse) WriteFile(path string) error {
	if data, err := os.ReadFile(path); err != nil {
		return err
	} else {
		return resp.response.Write(bytes.NewBuffer(data))
	}
}

func (resp *HttpResponse) WriteBinary(content []byte) error {
	return resp.response.Write(bytes.NewBuffer(content))
}

func (resp *HttpResponse) ReadNoCloser() []byte {
	r := resp.response
	if data, err := io.ReadAll(r.Body); err != nil {
		return nil
	} else {
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		return data
	}
}
