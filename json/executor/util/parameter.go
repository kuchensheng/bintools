package util

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/rs/zerolog/log"
	"github.com/yalp/jsonpath"
	"reflect"
	"strings"
)

func GetContextValue(ctx *gin.Context, key string) any {
	if strings.HasPrefix(key, consts.KEY_REQ) {
		if strings.HasPrefix(key, consts.KEY_REQ_BODY) {
			return GetBodyParameterValue(ctx, key)
		} else {
			return GetNotBodyParameterValue(ctx, key)
		}
	} else if !strings.HasPrefix(key, consts.KEY_TOKEN) {
		return key
	}
	return GetResultValue(ctx, key)
}

func GetBodyParameterValue(ctx *gin.Context, key string) any {
	if v, ok := ctx.Get(consts.PARAMETERMAP); !ok {
		return nil
	} else if body, existed := v.(map[string]any)[consts.KEY_REQ_BODY]; existed {
		if consts.KEY_REQ_BODY == key || key == "" {
			return body
		}
		split := strings.Split(key, consts.KEY_REQ_CONNECTOR)
		if result, ok1 := ReadByJsonPath(body.([]byte), split[2:]); ok1 {
			return result
		}
	}
	return nil
}

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

func SetResultValue(ctx *gin.Context, key string, value any) {
	if value == nil {
		log.Info().Msgf("key = %s,value is nil，不进行任何动作", key)
		return
	}
	//log.Info().Msgf("结果赋值,key=%s,value = %s", key, value)
	if v, ok := ctx.Get(consts.RESULTMAP); ok {
		data := value
		typeOf := reflect.TypeOf(value)
		log.Info().Msgf("value的类型:%v", typeOf)
		if typeOf != reflect.TypeOf([]byte("")) {
			data, _ = json.Marshal(value)
		}
		v.(map[string]any)[key] = data
	}
}

func ReadByJsonPath(v []byte, express []string) (any, bool) {
	var key []string
	key = append(key, "$")
	key = append(key, express...)
	path := strings.Join(key, ".")
	log.Info().Msgf("内容读取路径:key=%s", path)
	var data interface{}
	_ = json.Unmarshal(v, &data)
	if res, err := jsonpath.Read(data, path); err != nil {
		log.Warn().Msgf("无法从jsonPath中读取数据,%v", err)
		return nil, false
	} else {
		log.Info().Msgf("读取到内容:%v", res)
		return res, true
	}
}
