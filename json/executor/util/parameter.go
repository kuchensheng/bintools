package util

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/rs/zerolog/log"
	"github.com/yalp/jsonpath"
	"strings"
)

func GetContextValue(ctx *gin.Context, key string) any {
	if strings.HasPrefix(key, consts.KEY_REQ) {
		if strings.HasPrefix(key, consts.KEY_REQ_BODY) {
			return GetBodyParameterValue(ctx, key)
		} else {
			return GetNotBodyParameterValue(ctx, key)
		}
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
		if result, ok1 := ReadByJsonPath(body, split[2:]); ok1 {
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
		split := strings.Split(key, consts.KEY_REQ_CONNECTOR)
		mapKey := strings.Join(split[0:2], consts.KEY_REQ_CONNECTOR)
		if res, existed := resultMap[mapKey]; existed {
			if result, ok1 := ReadByJsonPath(res, split[2:]); ok1 {
				return result
			}
		}
	}
	return nil
}

func SetResultValue(ctx *gin.Context, key string, value any) {
	if v, ok := ctx.Get(consts.RESULTMAP); ok {
		v.(map[string]any)[key] = value
	}
}

func ReadByJsonPath(v any, express []string) (any, bool) {
	var key []string
	key = append(key, "$")
	key = append(key, express...)
	path := strings.Join(key, ".")
	log.Info().Msgf("内容读取路径:key=%s", path)
	if res, err := jsonpath.Read(v, path); err != nil {
		return nil, false
	} else {
		log.Info().Msgf("读取到内容:%v", res)
		return res, true
	}
}
