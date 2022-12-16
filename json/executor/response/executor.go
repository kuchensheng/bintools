package response

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/executor/util"
	"github.com/kuchensheng/bintools/json/model"
)

//BuildSuccessResponse 组装响应结果
func BuildSuccessResponse(ctx *gin.Context, responses map[string]model.ApixResponse) (any, error) {
	for s, response := range responses {
		if s == "200" {
			return readSchema(ctx, response.Schema), nil
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
