package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/kuchensheng/bintools/tracer/trace"
	"net/http"
	"net/url"
	"strings"
)

//执行服务节点
func ExecServer(ctx *gin.Context, step model.ApixStep) (any, error) {
	v, _ := ctx.Get(consts.TRACER)
	tracer := v.(*trace.ServerTracer)
	if request, err := buildRequest(ctx, step); err != nil {
		return nil, err
	} else {
		return tracer.Call(request)
	}
}

func buildRequest(ctx *gin.Context, step model.ApixStep) (*http.Request, error) {
	scheme := "http" //
	if step.Protocol == "https" {
		scheme = "https://"
	}
	domain := strings.ReplaceAll(step.Domain, "/", "")
	path := step.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	strUrl := fmt.Sprintf("%s%s%s", scheme, domain, path)
	url, _ := url.Parse(strUrl)
	request := &http.Request{
		Method: step.Method,
		URL:    url,
	}
	for _, parameter := range step.Parameters {
		location := parameter.In
		switch location {
		case consts.KEY_BODY:
			schema := parameter.Schema
			schemaType := schema.Type
			if schemaType == consts.OBJECT {
				body := make(map[string]any)
				for _, property := range schema.Properties {
					body[property.Name] = getValue(ctx, property.Default)
				}
				data, _ := json.Marshal(body)
				if r, e := http.NewRequest(step.Method, strUrl, bytes.NewBuffer(data)); e == nil {
					request = r
				} else {
					return nil, e
				}
			}
		case consts.KEY_QUERY:
			if v := getValue(ctx, parameter.Default); v != nil {
				url.Query().Add(parameter.Name, v.(string))
			}
		case consts.KEY_HEADER:
			if v := getValue(ctx, parameter.Default); v != nil {
				request.Header.Set(parameter.Name, v.(string))
			}
		case consts.KEY_COOKIE:
			if v := getValue(ctx, parameter.Default); v != nil {
				request.AddCookie(&http.Cookie{
					Name:  parameter.Name,
					Value: v.(string),
				})
			}
		case consts.KEY_FORM:
			if v := getValue(ctx, parameter.Default); v != nil {
				request.Form.Add(parameter.Name, v.(string))
			}
		default:
			return nil, errors.New("暂不支持的参数形式")
		}
	}
	return request, nil
}

func getValue(ctx *gin.Context, express string) any {
	//todo 从上下文中读取值
	return nil
}
