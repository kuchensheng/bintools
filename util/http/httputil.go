package http

import "strings"

const (
	//prefixHttps https请求前缀
	prefixHttps = "https:"

	//prefixHttp https请求前缀
	prefixHttp = "http:"
)

// IsHttps 检测是否https
func IsHttps(url string) bool {
	return strings.HasPrefix(strings.ToLower(url), prefixHttps)
}

// IsHttp 检测是否http
func IsHttp(url string) bool {
	return strings.HasPrefix(strings.ToLower(url), prefixHttp)
}

//CreateRequest 创建Http请求对象
