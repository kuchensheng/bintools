package util

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/rs/zerolog/log"
	"github.com/yalp/jsonpath"
	"reflect"
	"runtime/debug"
	"strings"
)

//GetContextValue 从上下文中读取参数或者结果值,如果为空则返回nil
func GetContextValue(ctx *gin.Context, key string) any {
	log.Info().Msgf("从当前上下文中获取值，key=[%s]", key)
	//判断是否以$req开头，true标识从请求参数中获取值
	if strings.HasPrefix(key, consts.KEY_REQ) {
		//判断是否以$req.data开头，true表示从请求体中获取值，否则从其他请求参数中获取参数
		if strings.HasPrefix(key, consts.KEY_REQ_BODY) {
			return GetBodyParameterValue(ctx, key)
		} else {
			return GetNotBodyParameterValue(ctx, key)
		}
	} else if !strings.HasPrefix(key, consts.KEY_TOKEN) {
		//如果key不是以$开头，不认识的常量或变量，则直接返回key
		return key
	}
	//从结果集中读取值
	return GetResultValue(ctx, key)
}

//GetBodyParameterValue 从请求体中参数值
func GetBodyParameterValue(ctx *gin.Context, key string) any {
	log.Info().Msgf("从请求体中获取值,key=[%s]", key)
	defer func() {
		if x := recover(); x != nil {
			log.Error().Msgf("读取请求体参数时异常,%v", x.(error))
			fmt.Printf("%s\n", debug.Stack())
		}
	}()
	if !strings.HasPrefix(key, consts.KEY_TOKEN) {
		return key
	}
	if v, ok := ctx.Get(consts.PARAMETERMAP); !ok {
		log.Warn().Msg("请求体不存在")
		return nil
	} else if body, existed := v.(map[string]any)[consts.KEY_REQ_BODY]; existed {
		if consts.KEY_REQ_BODY == key || key == "" {
			return body
		}
		split := strings.Split(key, consts.KEY_REQ_CONNECTOR)
		var express []string
		if len(split) > 2 {
			express = split[2:]
		}
		if result, ok1 := ReadByJsonPath(body.([]byte), express); ok1 {
			return result
		}
	}
	return nil
}

//GetNotBodyParameterValue 从非请求体的参数列表中获取值，如果不去不到，返回nil
func GetNotBodyParameterValue(ctx *gin.Context, key string) any {
	if key == "" {
		return nil
	}
	if v, ok := ctx.Get(consts.PARAMETERMAP); !ok {
		return nil
	} else if res, existed := v.(map[string]any)[key]; existed {
		return res
	}
	return nil
}

//GetResultValue 从各节点的结果集中读取值，key = $GraphId.$resp.export开头
func GetResultValue(ctx *gin.Context, key string) any {
	if v, ok := ctx.Get(consts.RESULTMAP); !ok {
		return nil
	} else {
		resultMap := v.(map[string]any)
		if result, find := resultMap[key]; find {
			return result
		}
		var suffix []string
		return getValue(resultMap, key, suffix)
	}
	return nil
}

func getValue(resultMap map[string]any, prefix string, suffix []string) any {
	if prefix == "" {
		return nil
	}
	if result, ok := resultMap[prefix]; ok {
		if suffix != nil {
			if result, ok = ReadByJsonPath(result.([]byte), suffix); ok {
				return result
			}
		}
		return result
	}
	split := strings.Split(prefix, consts.KEY_REQ_CONNECTOR)
	length := len(split)
	mapKey := strings.Join(split[0:length-1], consts.KEY_REQ_CONNECTOR)
	suffix = append(suffix, split[length-1])
	return getValue(resultMap, mapKey, suffix)
}

//SetResultValue 给当前上下文中的resultMap键赋值
func SetResultValue(ctx *gin.Context, key string, value any) {
	//if value == nil {
	//	log.Info().Msgf("key = %s,value is nil，不进行任何动作", key)
	//	return
	//}
	//log.Info().Msgf("结果赋值,key=%s,value = %s", key, value)
	if v, ok := ctx.Get(consts.RESULTMAP); ok {
		data := value
		typeOf := reflect.TypeOf(value)
		log.Info().Msgf("value的类型:%v", typeOf)
		if typeOf != reflect.TypeOf([]byte("")) {
			data, _ = json.Marshal(value)
		}
		v.(map[string]any)[key] = data
	} else {
		valueMap := make(map[string]any)
		valueMap[key] = value
		ctx.Set(consts.RESULTMAP, valueMap)
	}
}

//ReadByJsonPath 利用jsonPath读取对应的内容
func ReadByJsonPath(v []byte, express []string) (any, bool) {
	log.Info().Msgf("使用jsonPath解析值,%v", express)
	var key []string
	key = append(key, "$")
	key = append(key, express...)
	path := strings.Join(key, ".")
	log.Info().Msgf("内容读取路径:key=%s", path)
	var data interface{}
	_ = json.Unmarshal(v, &data)
	if res, err := jsonpath.Read(data, path); err != nil {
		//log.Warn().Msgf("无法从jsonPath中读取数据,%v", err)
		return nil, false
	} else {
		log.Info().Msgf("读取到内容:%v", res)
		return res, true
	}
}
