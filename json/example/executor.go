package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/kuchensheng/bintools/tracer/trace"
	"github.com/yalp/jsonpath"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var path = "/api/test/{1}"
var parameters = func() []model.ApixParameter {
	var parameters []model.ApixParameter
	//将字符串内容初始化
	var strParameter = "[]"
	json.Unmarshal([]byte(strParameter), &parameters)
	return parameters
}()

var response = func() map[string]model.ApixResponse {
	responseMap := make(map[string]model.ApixResponse)
	//字符串初始化
	var strResponse = "{}"
	json.Unmarshal([]byte(strResponse), &responseMap)
	return responseMap
}()

var steps = func() []model.ApixStep {
	var steps []model.ApixStep
	//将字符串内容初始化
	var strStep = "[]"
	json.Unmarshal([]byte(strStep), &steps)
	return steps
}()

var scriptEngine *goja.Runtime

var resultMap = make(map[string][]byte)
var parameterMap = make(map[string][]byte)

//Executor 插件的执行入口
//export Executor
func Executor(context *gin.Context) {
	//参数检查
	if exception := checkParameter(context); exception != nil {
		context.JSON(400, exception)
		return
	}
	//开启tracer
	serverTracer := trace.NewServerTracer(context.Request)
	exception := executeStep(steps2Map(steps), serverTracer)
	if exception != nil {
		context.JSON(400, exception)
	}
	//组装返回值
}

func steps2Map(steps []model.ApixStep) map[string]model.ApixStep {
	result := make(map[string]model.ApixStep)
	for _, step := range steps {
		result[step.GraphId] = step
	}
	return result
}

func executeStep(steps map[string]model.ApixStep, tracer *trace.ServerTracer) *model.BusinessException {
	var exception *model.BusinessException
	var roots []model.ApixStep
	for _, step := range steps {
		if step.PrevId == "" {
			roots = append(roots, step)
		}
	}
	ch := make(chan *model.BusinessException, len(roots))
	for _, root := range roots {
		go func(t *trace.ServerTracer, channel chan *model.BusinessException) {
			channel <- executeCurrentStep(root, steps, t)
		}(tracer, ch)
	}
	for i := 0; i < len(roots); i++ {
		select {
		case exception = <-ch:
		case <-time.After(5 * time.Second):
			exception = model.NewBusinessException(1080500, "执行超时了")
		}
	}
	return exception
}

//executePredicate 筛选出执行步骤
func executeCurrentStep(step model.ApixStep, steps map[string]model.ApixStep, tracer *trace.ServerTracer) *model.BusinessException {
	if stepIsEmpty(step) {
		return nil
	}

	predicates := step.Predicate
	//下个节点
	next := steps[step.ThenGraphId]
	//执行当前节点
	if !scriptIsEmpty(step.Script) {
		if data, err := executeGoScript(step.Script); err != nil {
			return model.NewBusinessException(1080500, "脚本节点执行失败:"+err.Error())
		} else {
			resultMap[step.GraphId] = data
		}
	} else if predicates != nil && len(predicates) > 0 {
		next = steps[findNexStep(step)]
	} else {
		//服务：执行当前节点，发起http请求
		if exception := executeServerStep(step, tracer); exception != nil {
			return exception
		}
	}
	return executeCurrentStep(next, steps, tracer)
}

func executeServerStep(step model.ApixStep, tracer *trace.ServerTracer) *model.BusinessException {
	proto := step.Protocol
	if proto == "" {
		proto = "http"
	}
	if strings.Contains(proto, "://") {
		proto = strings.ReplaceAll(proto, "://", "")
	}

	requestUrl, err := url.Parse(fmt.Sprintf("%s://%s%s", proto, step.Domain, fillPath(step.Path, step.Parameters)))
	if err != nil {
		return model.NewBusinessException(1080500, "url拼写不正确")
	}

	request := &http.Request{
		Method: step.Method,
		URL:    requestUrl,
	}
	fillHeader(step.Parameters, request)
	fillForm(step.Parameters, request)
	fillQuery(step.Parameters, request)
	fillCookie(step.Parameters, request)
	//todo 发起调用
	return nil
}

func fillBody(parameters []model.ApixParameter, request *http.Request) {
	if parameters != nil {
		for _, parameter := range parameters {
			if parameter.In == "body" {
				schema := parameter.Schema
				if schema.Type == "object" {

				}
				break
			}
		}
	}
}

func parse2Body() {

}

func fillCookie(parameters []model.ApixParameter, request *http.Request) {
	if parameters != nil {
		for _, parameter := range parameters {
			if parameter.In == "cookie" {
				request.AddCookie(&http.Cookie{
					Name:  parameter.Name,
					Value: getValueByKey(parameter.Default),
				})
			}
		}
	}
}

func fillQuery(parameters []model.ApixParameter, request *http.Request) {
	if parameters != nil {
		for _, parameter := range parameters {
			if parameter.In == "query" {
				request.URL.Query().Add(parameter.Name, getValueByKey(parameter.Default))
			}
		}
	}
}

func fillForm(parameters []model.ApixParameter, request *http.Request) {
	if parameters != nil {
		for _, parameter := range parameters {
			if parameter.In == "formData" {
				request.Form.Add(parameter.Name, getValueByKey(parameter.Default))
			}
		}
	}
}

func fillHeader(parameters []model.ApixParameter, request *http.Request) {
	if parameters != nil {
		for _, parameter := range parameters {
			if parameter.In == "header" {
				request.Header.Set(parameter.Name, getValueByKey(parameter.Default))
			}
		}
	}
}

func fillPath(path string, parameters []model.ApixParameter) string {
	//根据参数定义填充url
	if parameters != nil {
		for _, parameter := range parameters {
			if parameter.In == "path" {
				path = strings.Replace(path, parameter.Name, getValueByKey(parameter.Default), 0)
			}
		}
	}
	return path
}

func findNexStep(step model.ApixStep) string {
	//predicateType = 0,表示所有条件都要满足，1=任一条件为真
	predicateType := step.PredicateType
	predicates := step.Predicate
	predicateValue := len(predicates)

	for _, predicate := range predicates {
		if !predicate.Enabled {
			continue
		}
		if predicate.Type == "if" {
			//todo 这里是操作符拼接
			if getValueByKey(predicate.Key) != getValueByKey(predicate.Value) {
				predicateValue -= 1
			}
		} else {
			thenGraphId := func(cases []model.ApixSwitchPredicate) string {
				for _, switchPredicate := range cases {
					//todo 这里是操作符拼接
					if getValueByKey(switchPredicate.Key) == getValueByKey(switchPredicate.Value) {
						return switchPredicate.ThenGraphId
					}
				}
				return ""
			}(predicate.Cases)
			if thenGraphId != "" {
				step.ThenGraphId = thenGraphId
			} else {
				predicateValue -= 1
			}
		}
	}
	if predicateType == 0 && predicateValue < len(predicates) {
		return step.ElseGraphId
	}
	return step.ThenGraphId
}

func getValueByKey(key string) string {
	if !(strings.Contains(key, ".") || strings.Contains(key, "$")) {
		return key
	}
	key = strings.ReplaceAll(key, "#", ".")
	splits := strings.Split(key, ".")
	if len(splits) < 3 {
		//key值错误，返回空
		return key
	}
	graphId := splits[0]
	valueMap := make(map[string][]byte)
	location := splits[1]
	if location == "$resp" {
		valueMap = resultMap
	} else {
		valueMap = parameterMap
	}
	subKey := strings.Join(splits[2:], ".")
	data := valueMap[graphId]
	if data == nil {
		//值不存在
		return ""
	} else if v, err := jsonpath.Read(data, subKey); err != nil {
		println(err)
		return ""
	} else {
		return v.(string)
	}
}

func executeGoScript(script model.ApixScript) ([]byte, error) {
	defer deferHandler()
	if scriptEngine == nil {
		scriptEngine = goja.New()
	}
	v, err := scriptEngine.RunString(script.Script)
	if err != nil {
		return nil, err
	}
	return []byte(v.String()), nil
}

func deferHandler() error {
	if x := recover(); x != nil {
		return x.(error)
	}
	return nil
}

func predicateIsEmpty(predicates []model.ApiStepPredicate) bool {
	return predicates == nil || len(predicates) == 0
}

func scriptIsEmpty(script model.ApixScript) bool {
	return reflect.DeepEqual(script, model.ApixScript{})
}

func stepIsEmpty(step model.ApixStep) bool {
	return reflect.DeepEqual(step, model.ApixStep{})
}

//todo 暂时不需要进行参数检查
func checkParameter(context *gin.Context) *model.BusinessException {
	defer deferHandler()
	for _, parameter := range parameters {
		location := parameter.In
		parameterName := parameter.Name
		required := parameter.Required
		parameterType := parameter.Type
		switch location {
		case "header":
			//从请求头中获取
			return checkHeader(context, parameterName, required)
		case "body":
			//读取请求体内容
			data, exception := readRequestBody(context)
			if exception != nil {
				return exception
			}
			schema := parameter.Schema
			subType := schema.SubType
			switch schema.Type {
			case "object":
				//object类型，直接校验字段信息
				return checkProperty(schema.Properties, parameterName, parameterType, required, data)
			case "array":
				if subType == "object" {
					for _, child := range schema.Children {
						//解析child，child是个对象
						if exception := checkProperty(child.Properties, child.Name, child.Type, child.Required, nil); exception != nil {
							return exception
						}
					}
				} else {
					//subType是基本类型，直接读取元素内容
				}
			default:
				//基本类型

			}

		case "cookie":
			return checkCookie(context, parameterName, required)
		case "formData":
			return checkFormData(context, parameterName, parameterType, required)
		case "query":
			return checkQuery(context, parameterName, parameterType, required)
		default:
			//未知的参数定义
			return nil
		}
	}
	return nil
}

func readRequestBody(context *gin.Context) ([]byte, *model.BusinessException) {
	defer deferHandler()
	data, err := context.GetRawData()
	if err != nil {
		return nil, model.NewBusinessException(1080500, "无法读取请求体")
	}
	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return data, nil
}

func checkQuery(context *gin.Context, parameterName, parameterType string, required bool) *model.BusinessException {
	defer deferHandler()
	switch parameterType {
	case "array":
		if a := context.QueryArray(parameterName); (a == nil || len(a) == 0) && required {
			return model.NewBusinessException(1080400, "必填参数缺失:["+parameterName+"]")
		}
	case "object":
		fallthrough
	case "map":
		if obj := context.QueryMap(parameterName); (obj == nil || len(obj) == 0) && required {
			return model.NewBusinessException(1080400, "必填参数缺失:["+parameterName+"]")
		}
	default:
		if v := context.Query(parameterName); v == "" && required {
			return model.NewBusinessException(1080400, "必填参数缺失:["+parameterName+"]")
		}
	}
	return nil
}

func checkHeader(context *gin.Context, parameterName string, required bool) *model.BusinessException {
	defer deferHandler()
	header := context.GetHeader(parameterName)
	if header == "" && required {
		return model.NewBusinessException(1080400, "必填参数缺失:["+parameterName+"]")
	}
	return nil
}

func checkCookie(context *gin.Context, parameterName string, required bool) *model.BusinessException {
	defer deferHandler()
	if c, err := context.Request.Cookie(parameterName); err != nil {
		return model.NewBusinessException(1080500, "无法从cookie中读取参数信息:"+err.Error())
	} else if c == nil && required {
		return model.NewBusinessException(1080400, "必填参数缺失cookie:["+parameterName+"]:"+err.Error())
	}
	return nil
}

func checkFormData(context *gin.Context, parameterName, parameterType string, required bool) *model.BusinessException {
	defer deferHandler()
	switch parameterType {
	case "array":
		if arrayP := context.PostFormArray(parameterName); arrayP == nil || len(arrayP) == 0 {
			if required {
				return model.NewBusinessException(1080400, "必填参数缺失:["+parameterName+"]")
			}
		}
	case "object":
		if objP := context.PostFormMap(parameterName); objP == nil {
			if required {
				return model.NewBusinessException(1080400, "必填参数缺失:["+parameterName+"]")
			}
		}
	default:
		if v := context.PostForm(parameterName); v == "" && required {
			return model.NewBusinessException(1080400, "必填参数缺失:["+parameterName+"]")
		}
	}
	return nil
}

//checkProperty 检查属性信息，主要检查是否遗漏必填参数
func checkProperty(properties map[string]model.ApixProperty, parameterName, parameterType string, required bool, data []byte) *model.BusinessException {
	defer deferHandler()
	if required && (nil == properties || len(properties) == 0) {
		return model.NewBusinessException(1080400, "请求体不能为空")
	}
	if data == nil && !required {
		//todo 无内容
		return nil
	}
	bodyValue := make(map[string][]byte)
	if err := json.Unmarshal(data, &bodyValue); err != nil {
		return model.NewBusinessException(1080500, "无法解析请求体内容："+err.Error())
	}
	//paramter存在吗？
	existed := false
	for s, property := range properties {
		if s == parameterName {
			existed = true
			break
		}

		switch property.Type {
		case "object":
			if exception := checkProperty(properties, s, property.Type, property.Required, bodyValue[property.Name]); exception != nil {
				return exception
			}
		case "array":
			//todo 暂不实现
		default:
			//读取属性值

		}

	}
	if !existed && required {
		return model.NewBusinessException(1080400, "必填参数缺失:["+parameterName+"]")
	}
	return nil
}
