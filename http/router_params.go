package http

import (
	"mime/multipart"
	"net/http"
	"net/url"
)

// 如何传a int,b string,c int32,d other and so on 这样的不定长参数？
type HandlerParamFunc func(params []HandlerParam) (any, error)

type HandlerParam interface {
	Name() string
	Required() bool
	Value() any
}

type HeaderParam struct {
	name     string
	value    string
	required bool
}

func (p HeaderParam) Name() string {
	return p.name
}

func (p HeaderParam) Value() any {
	return p.value
}

func (p HeaderParam) Required() bool {
	return true
}

type PathParam struct {
	name     string
	value    any
	required bool
}

func (p PathParam) Name() string {
	return p.name
}

func (p PathParam) Value() any {
	return p.value
}

func (p PathParam) Required() bool {
	return true
}

type QueryParam struct {
	name     string
	value    any
	required bool `json:"required,omitempty"`
}

func NewQuery(name string, required bool) QueryParam {
	return QueryParam{
		name:     name,
		required: required,
	}
}

func (p QueryParam) Name() string {
	return p.name
}

func (p QueryParam) Value() any {
	return p.value
}

func (p QueryParam) Required() bool {
	return p.required
}

type BodyParam struct {
	Body     any
	required bool `json:"required,omitempty"`
}

func (b BodyParam) Name() string {
	return "body"
}

func (b BodyParam) Required() bool {
	return b.required
}

func (b BodyParam) Value() any {
	return b.Body
}

type FormParam struct {
	Form     url.Values
	required bool `json:"required,omitempty"`
}

func (b FormParam) Name() string {
	return "form"
}

func (b FormParam) Required() bool {
	return b.required
}

func (b FormParam) Value() any {
	return b.Form
}

type MultiFormParam struct {
	Form     *multipart.Form
	required bool `json:"required,omitempty"`
}

func (b MultiFormParam) Name() string {
	return "multipart-form"
}

func (b MultiFormParam) Required() bool {
	return b.required
}

func (b MultiFormParam) Value() any {
	return b.Form
}

func (e *engine) GetWithParam(pattern string, paramFunc HandlerParamFunc) {
	e.registerRouterWithParams(http.MethodGet, pattern, paramFunc)
}

func (e *engine) PostWithParam(pattern string, paramFunc HandlerParamFunc) {
	e.registerRouterWithParams(http.MethodPost, pattern, paramFunc)
}

func (e *engine) PutWithParam(pattern string, paramFunc HandlerParamFunc) {
	e.registerRouterWithParams(http.MethodPut, pattern, paramFunc)
}

func (e *engine) DeleteWithParam(pattern string, paramFunc HandlerParamFunc) {
	e.registerRouterWithParams(http.MethodDelete, pattern, paramFunc)
}

func (e *engine) OptionsWithParam(pattern string, paramFunc HandlerParamFunc) {
	e.registerRouterWithParams(http.MethodOptions, pattern, paramFunc)
}

func (e *engine) HeadWithParam(pattern string, paramFunc HandlerParamFunc) {
	e.registerRouterWithParams(http.MethodHead, pattern, paramFunc)
}

func (e *engine) AnyWithParam(pattern string, paramFunc HandlerParamFunc) {
	e.registerRouterWithParams(http.MethodGet, pattern, paramFunc)
	e.registerRouterWithParams(http.MethodPost, pattern, paramFunc)
	e.registerRouterWithParams(http.MethodPut, pattern, paramFunc)
	e.registerRouterWithParams(http.MethodDelete, pattern, paramFunc)
	e.registerRouterWithParams(http.MethodOptions, pattern, paramFunc)
	e.registerRouterWithParams(http.MethodHead, pattern, paramFunc)
}
