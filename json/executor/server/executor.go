package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/executor/util"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/kuchensheng/bintools/tracer/trace"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//执行服务节点
func ExecServer(ctx *gin.Context, step model.ApixStep) (any, error) {
	v, _ := ctx.Get(consts.TRACER)
	tracer := v.(*trace.ServerTracer)
	if request, err := buildRequest(ctx, step); err != nil {
		log.Warn().Msgf("不能正确地构建请求")
		return nil, consts.NewException(step.GraphId, "", err.Error())
	} else if request != nil {
		log.Info().Msgf("请求地址:%s", request.URL.String())
		if result, err1 := tracer.Call(request); err1 != nil {
			log.Warn().Msgf("服务节点执行失败,%v", err1)
			return nil, consts.NewException(step.GraphId, "", err1.Error())
		} else {
			return result, nil
		}
	}
	return nil, nil
}

func buildRequest(ctx *gin.Context, step model.ApixStep) (*http.Request, error) {
	scheme := "http://" //
	if step.Protocol == "https" {
		scheme = "https://"
	}
	if step.Path == "" || step.Domain == "" || step.Method == "" {
		return nil, nil
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
		Header: make(map[string][]string),
		Form:   make(map[string][]string),
	}
	for _, parameter := range step.Parameters {
		location := parameter.In
		switch location {
		case consts.KEY_BODY:
			schema := parameter.Schema
			schemaType := schema.Type
			if schemaType == consts.OBJECT && schema.Properties != nil && len(schema.Properties) > 0 {
				body := make(map[string]any)
				for _, property := range schema.Properties {
					if v := util.GetBodyParameterValue(ctx, property.Default); v != nil {
						body[property.Name] = util.GetBodyParameterValue(ctx, property.Default)
					}
				}
				data, _ := json.Marshal(body)
				request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				if len(body) == 0 {
					request.ContentLength = 0
				}
			}

		case consts.KEY_QUERY:
			if v := util.GetNotBodyParameterValue(ctx, parameter.Default); v != nil {
				url.Query().Add(parameter.Name, v.(string))
			}
		case consts.KEY_HEADER:
			if v := util.GetNotBodyParameterValue(ctx, parameter.Default); v != nil {
				request.Header.Set(parameter.Name, v.(string))
			}
		case consts.KEY_COOKIE:
			if v := util.GetNotBodyParameterValue(ctx, parameter.Default); v != nil {
				request.AddCookie(&http.Cookie{
					Name:  parameter.Name,
					Value: v.(string),
				})
			}
		case consts.KEY_FORM:
			if v := util.GetNotBodyParameterValue(ctx, parameter.Default); v != nil {
				request.Form.Add(parameter.Name, v.(string))
			}
		default:
			return nil, errors.New("暂不支持的参数形式")
		}
	}
	return request, nil
}
