//go:build windows

package example

import "C"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/kuchensheng/bintools/tracer/trace"
	"github.com/yalp/jsonpath"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
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

// Executor 插件的执行入口
//
//export Executor
func Executor(request *C.int, response *C.int) {
	r := (*http.Request)(unsafe.Pointer(request))
	w := *(*http.Response)(unsafe.Pointer(response))
	//参数检查
	if exception := checkParameter(context); exception != nil {
		context.JSON(400, exception)
		return
	}
	//开启tracer
	serverTracer := trace.NewServerTracer(context.Request)
	defer serverTracer.EndTraceOk()
	exception := executeStep(steps2Map(steps), serverTracer)
	if exception != nil {
		context.JSON(400, exception)
		return
	}
	//组装返回值
	result := packingResponse()
	context.JSON(200, result)
}

func packingResponse() map[string]any {
	for _, apixResponse := range response {
		schema := apixResponse.Schema
		return parseProperties2Body(schema.Properties)
	}
	return nil
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

// executePredicate 筛选出执行步骤
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
	//todo 组装请求体
	fillBody(step.Parameters, request)
	if data, err := tracer.Call(request); err != nil {
		return model.NewBusinessException(1080500, "请求调用失败:"+err.Error())
	} else {
		resultMap[step.GraphId] = data
	}
	return nil
}

func fillBody(parameters []model.ApixParameter, r *http.Request) {
	if parameters != nil {
		for _, parameter := range parameters {
			if parameter.In == "body" {
				schema := parameter.Schema
				var body any
				switch schema.Type {
				case "object":
					body = parseProperties2Body(schema.Properties)
				case "array":
					body = parseProperty2Body(schema.Children)
				default:
					body = getValueByKey(schema.Default)
				}
				data, _ := json.Marshal(body)
				r.GetBody = func() (io.ReadCloser, error) {
					return ioutil.NopCloser(bytes.NewBuffer(data)), nil
				}
				break
			}
		}
	}
}

func parseProperty2Body(properties []model.ApixProperty) []any {
	var childList []any
	for _, child := range properties {
		switch child.Type {
		case "object":
			childList = append(childList, parseProperties2Body(child.Properties))
		case "array":
			childList = append(childList, parseProperty2Body(child.Children))
		default:
			childList = append(childList, getValueByKey(child.Default))
		}
	}
	return childList
}

func parseProperties2Body(properties map[string]model.ApixProperty) map[string]any {
	body := make(map[string]any)
	for name, property := range properties {
		switch property.Type {
		case "object":
			body[name] = parseProperties2Body(property.Properties)
		case "array":
			//读取children
			body[name] = parseProperty2Body(property.Children)
		default:
			body[name] = getValueByKey(property.Default)
		}
	}
	return body
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
			if checkKV(getValueByKey(predicate.Key), getValueByKey(predicate.Value), predicate.Operator) {
				predicateValue -= 1
			}
		} else {
			thenGraphId := func(cases []model.ApixSwitchPredicate) string {
				defaultCase := ""
				for _, switchPredicate := range cases {
					if switchPredicate.IsDefault {
						defaultCase = switchPredicate.ThenGraphId
					}
					if checkKV(getValueByKey(switchPredicate.Key), getValueByKey(switchPredicate.Value), switchPredicate.Operator) {
						return switchPredicate.ThenGraphId
					}
				}
				return defaultCase
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

func checkKV(key, value, operator string) bool {
	switch operator {
	case "==":
		return key == value
	case "!=":
		return key != value
	case ">":
		keyInt, valueInt := convertStr2Int(key, value)
		return keyInt > valueInt
	case ">=":
		keyInt, valueInt := convertStr2Int(key, value)
		return keyInt >= valueInt
	case "<":
		keyInt, valueInt := convertStr2Int(key, value)
		return keyInt < valueInt
	case "<=":
		keyInt, valueInt := convertStr2Int(key, value)
		return keyInt <= valueInt
	case "contains":
		return strings.Contains(key, value)
	case "not contains":
		return !strings.Contains(key, value)
	default:
		return true
	}
}

func convertStr2Int(key, value string) (keyInt, valueInt int) {
	keyInt, err := strconv.Atoi(key)
	if err != nil {
		log.Print("无法将key=[", key, "]转换为int类型")
		keyInt = 0
	}
	valueInt, err = strconv.Atoi(value)
	if err != nil {
		log.Print("无法将key=[", value, "]转换为int类型")
		valueInt = 0
	}
	return
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

// todo 暂时不需要进行参数检查
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
	parameterMap["data"] = data
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

// checkProperty 检查属性信息，主要检查是否遗漏必填参数
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

func main() {

}
