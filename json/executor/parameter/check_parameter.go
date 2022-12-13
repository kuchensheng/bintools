package parameter

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/model"
	"io"
	"io/ioutil"
)

func deferHandler() error {
	if x := recover(); x != nil {
		return x.(error)
	}
	return nil
}

func CheckParameter(ctx *gin.Context, parameters []model.ApixParameter) error {
	for _, parameter := range parameters {
		location := parameter.In
		switch location {
		case consts.KEY_QUERY:
			if e := checkQuery(ctx, parameter.Name, parameter.Required); e != nil {
				return e
			}
		case consts.KEY_HEADER:
			if e := checkHeader(ctx, parameter.Name, parameter.Required); e != nil {
				return e
			}
		case consts.KEY_FORM:
			if e := checkFormData(ctx, parameter.Name, parameter.Required); e != nil {
				return e
			}
		case consts.KEY_COOKIE:
			if e := checkCookie(ctx, parameter.Name, parameter.Required); e != nil {
				return e
			}
		case consts.KEY_BODY:
			if e := checkBody(ctx, parameter.Name, parameter.Required); e != nil {
				return e
			}
		default:
			return errors.New("不支持的类型")
		}
	}
	return nil
}

func readRequestBody(ctx *gin.Context) ([]byte, error) {
	defer deferHandler()
	r := ctx.Request
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, newError(consts.KEY_BODY, "无法读取请求体")
	}
	r.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewBuffer(data)), nil
	}
	return data, nil
}

func checkBody(ctx *gin.Context, parameterName string, required bool) error {
	defer deferHandler()
	if body, e := readRequestBody(ctx); e != nil {
		return e
	} else if v, ok := ctx.Get(consts.PARAMETERMAP); ok {
		v.(map[string]any)[consts.KEY_REQ_BODY] = body
	}
	return nil
}

func checkQuery(ctx *gin.Context, parameterName string, required bool) error {
	defer deferHandler()
	r := ctx.Request
	get := r.URL.Query().Get(parameterName)
	if get == "" && required {
		return newError(consts.KEY_QUERY, parameterName)
	}
	if v, ok := ctx.Get(consts.PARAMETERMAP); ok {
		v.(map[string]any)[consts.KEY_REQ_QUERY+consts.KEY_REQ_CONNECTOR+parameterName] = get
	}
	return nil
}

func newError(location, name string) error {
	return errors.New(fmt.Sprintf("%s参数缺失，%s=null", location, name))
}

func checkFormData(ctx *gin.Context, parameterName string, required bool) error {
	defer deferHandler()
	r := ctx.Request
	get := r.Form.Get(parameterName)
	if get == "" {
		v := r.MultipartForm.Value
		if v1, ok := ctx.Get(consts.PARAMETERMAP); ok {
			v1.(map[string]any)[consts.KEY_REQ_QUERY+consts.KEY_REQ_CONNECTOR+parameterName] = v
		}
		if v == nil && required {
			return newError(consts.KEY_FORM, parameterName)
		}
	}
	if v1, ok := ctx.Get(consts.PARAMETERMAP); ok {
		v1.(map[string]any)[consts.KEY_REQ_QUERY+consts.KEY_REQ_CONNECTOR+parameterName] = get
	}
	return nil
}

func checkHeader(ctx *gin.Context, parameterName string, required bool) error {
	defer deferHandler()
	request := ctx.Request
	header := request.Header.Get(parameterName)
	if header == "" && required {
		return newError("请求头", parameterName)
	}
	if v1, ok := ctx.Get(consts.PARAMETERMAP); ok {
		v1.(map[string]any)[consts.KEY_REQ_QUERY+consts.KEY_REQ_CONNECTOR+parameterName] = header
	}
	return nil
}

func checkCookie(ctx *gin.Context, parameterName string, required bool) error {
	defer deferHandler()
	r := ctx.Request
	if c, err := r.Cookie(parameterName); err != nil {
		return newError(consts.KEY_COOKIE, parameterName)
	} else if c == nil && required {
		return newError(consts.KEY_COOKIE, parameterName)
	} else {
		if v1, ok := ctx.Get(consts.PARAMETERMAP); ok {
			v1.(map[string]any)[consts.KEY_REQ_QUERY+consts.KEY_REQ_CONNECTOR+parameterName] = c
		}
	}
	return nil
}
