package http

import (
	"log"
	"net/http"
	"net/url"
)

type HttpRequest struct {
	http.Request
}

func CreateRequest(method, rawURL string) *HttpRequest {
	return &HttpRequest{
		http.Request{
			Method: method,
			URL:    parseUrl(rawURL),
		},
	}
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
		log.Panicf("无法解析的url：%s", rawURL)
	}
	return parse
}
