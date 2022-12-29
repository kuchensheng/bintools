package response

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	log2 "github.com/kuchensheng/bintools/json/executor/log"
	"github.com/kuchensheng/bintools/json/executor/util"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/kuchensheng/bintools/tracer/trace"
	"github.com/rs/zerolog/log"
)

//BuildSuccessResponse 组装响应结果
func BuildSuccessResponse(ctx *gin.Context, responses map[string]model.ApixResponse) (any, error) {
	v, _ := ctx.Get(consts.TRACER)
	tracer := v.(*trace.ServerTracer)
	pk := log2.GetPackage(ctx)
	ls := log2.LogStruct{PK: pk, TraceId: tracer.TracId}
	ls.Info("开始组装响应结果...")
	defer func() {
		if x := recover(); x != nil {
			log.Warn().Msgf("结果组装异常:%v", x)
			ls.Error("结果组装失败:%s", x.(error).Error())
		} else {
			ls.Info("结果组装完毕")
		}
	}()
	for s, response := range responses {
		if s == "200" {
			schema := readSchema(ctx, response.Schema)
			log.Info().Msgf("组装结果:%s", schema)
			return schema, nil
		}
	}
	return nil, nil
}

func readSchema(ctx *gin.Context, schema model.ApixSchema) any {
	schemaType := schema.Type
	switch schemaType {
	case consts.OBJECT:
		result := make(map[string]any)
		for s, property := range schema.Properties {
			result[s] = readProperty(ctx, property)
		}
		return result
	case consts.ARRAY:
		var result []any
		for _, child := range schema.Children {
			result = append(result, readProperty(ctx, child))
		}
		return result
	default:
		return util.GetContextValue(ctx, schema.Default)
	}
}

func readProperty(ctx *gin.Context, property model.ApixProperty) any {
	propertyType := property.Type
	switch propertyType {
	case consts.OBJECT:
		result := make(map[string]any)
		for s, apixProperty := range property.Properties {
			result[s] = readProperty(ctx, apixProperty)
		}
		return result
	case consts.ARRAY:
		var result []any
		for _, child := range property.Children {
			result = append(result, readProperty(ctx, child))
		}
		return result
	default:
		return util.GetContextValue(ctx, property.Default)
	}
}
