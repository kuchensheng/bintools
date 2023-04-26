package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/kuchensheng/bintools/logger"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"time"
)

type HttpRequest struct {
	request *http.Request
	client  *http.Client
}

var (
	log = logger.GlobalLogger
)

func CreateRequest(method, rawURL string) *HttpRequest {
	return createHttpRequest(method, rawURL)
}

func CreateGet(rawURL string) *HttpRequest {
	return CreateRequest(http.MethodGet, rawURL)
}

func CreatePost(rawUrl string) *HttpRequest {
	return CreateRequest(http.MethodPost, rawUrl)
}

func parseUrl(rawURL string) *url.URL {
	parse, err := url.Parse(rawURL)
	if err != nil {
		log.Panic(fmt.Sprintf("无法解析的url：%s", rawURL))
	}
	return parse
}

// Get 发送Get请求，返回请求其内容字符串,异常返回nil
func Get(rawURL string) string {
	return GetWithTimeout(rawURL, timeout)
}

// GetWithTimeout 发送Get请求，设置请求超时时间。
// 返回请求其内容字符串,异常返回nil
func GetWithTimeout(rawURL string, timeout time.Duration) string {
	return CreateGet(rawURL).Timeout(timeout).Execute().GetString()
}

// GetWithParams 发送Get请求，设置请求参数。
// 返回请求其内容字符串,异常返回nil
func GetWithParams(rawURL string, params map[string]string) string {
	get := CreateGet(rawURL)
	return get.Params(params).Execute().GetString()
}

// Post 发送Post请求，返回请求其内容字符串,异常返回nil
func Post(rawURL, body string) string {
	return PostWithTimeout(rawURL, body, timeout)
}

// PostWithTimeout 发送Post请求，并设置超时时间 返回请求其内容字符串,异常返回nil
func PostWithTimeout(rawURL, body string, timeout time.Duration) string {
	return CreatePost(rawURL).Timeout(timeout).Body(body).Execute().GetString()
}

// PostWithForm 发送Post表单请求，并设置超时时间 返回请求其内容字符串,异常返回nil
func PostWithForm(rawURL, body string, form map[string]string) string {
	return CreatePost(rawURL).Form(form).Body(body).Execute().GetString()
}

func CreateDelete(rawURL string) *HttpRequest {
	return createHttpRequest(http.MethodDelete, rawURL)
}
func Delete(rawURL string) string {
	return CreateDelete(rawURL).Execute().GetString()
}

func DeleteWithParam(rawURL string, params map[string]string) string {
	return CreateDelete(rawURL).Params(params).Execute().GetString()
}

func Put(rawURL, body string) string {
	return createHttpRequest(http.MethodPut, rawURL).Body(body).Execute().GetString()
}

func PutWithParams(rawURL, body string, params map[string]string) string {
	return createHttpRequest(http.MethodPut, rawURL).Params(params).Body(body).Execute().GetString()
}

// Execute 执行请求
// 返回响应内容
func (request *HttpRequest) Execute() *HttpResponse {
	if do, err := request.client.Do(request.request); err != nil {
		logger.GlobalLogger.Error("request请求失败,method = %s,url=%s", request.request.Method, request.request.URL)
		return nil
	} else {
		return &HttpResponse{do}
	}
}

func createHttpClient(isHttps bool) *http.Client {
	client := http.DefaultClient
	if isHttps {
		client.Transport = &http.Transport{
			DialContext: func(dialer *net.Dialer) func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext
			}(&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}),
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	client.Timeout = timeout
	return client
}

func createHttpRequest(method, rawURL string) *HttpRequest {
	return &HttpRequest{
		&http.Request{
			Method:   method,
			URL:      parseUrl(rawURL),
			Header:   make(http.Header),
			Form:     make(url.Values),
			PostForm: make(url.Values),
		}, createHttpClient(IsHttps(rawURL)),
	}
}

func (request *HttpRequest) Timeout(customTimeout time.Duration) *HttpRequest {
	request.client.Timeout = customTimeout
	return request
}

func (request *HttpRequest) Head(header map[string]string) *HttpRequest {
	r := request.request
	for key, val := range header {
		r.Header.Set(key, val)
	}
	return request
}

func (request *HttpRequest) ContentType(contentType string) *HttpRequest {
	request.request.Header.Set("Content-Type", contentType)
	return request
}

func (request *HttpRequest) Params(params map[string]string) *HttpRequest {
	u := request.request.URL
	for key, val := range params {
		u.Query().Set(key, val)
	}
	return request
}

func (request *HttpRequest) Form(form map[string]string) *HttpRequest {
	f := request.request.Form
	if http.MethodGet != request.request.Method {
		f = request.request.PostForm
	}
	for key, val := range form {
		f.Set(key, val)
	}
	if http.MethodGet != request.request.Method {
		request.request.PostForm = f
	} else {
		request.request.Form = f
	}
	return request
}

func (request *HttpRequest) MultipartForm(form *multipart.Form) *HttpRequest {
	request.request.MultipartForm = form
	return request
}

func (request *HttpRequest) BodyBinary(body []byte) *HttpRequest {
	r := request.request
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return request
}

func (request *HttpRequest) ReadNoCloser() []byte {
	r := request.request
	if data, err := io.ReadAll(r.Body); err != nil {
		return nil
	} else {
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		return data
	}
}

func (request *HttpRequest) Body(body string) *HttpRequest {
	r := request.request
	r.Body = io.NopCloser(bytes.NewBufferString(body))
	return request
}
