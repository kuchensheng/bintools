package predicate

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/rs/zerolog/log"
)

const (
	eq   string = "=="
	ne   string = "!="
	ge   string = ">="
	gt   string = ">"
	le   string = "<="
	lt   string = "<"
	inc  string = "inc"
	_inc string = "!inc"
)

//ExecPredicates predicateType = 1 所有条件都为真，否则任一条件为真
func ExecPredicates(ctx *gin.Context, predicates []model.ApiStepPredicate, predicateType int) (bool, error) {
	log.Info().Msgf("执行逻辑判断，type= %d", predicateType)
	var b bool
	var paramMap = make(map[string]any)
	var resultMap = make(map[string]any)
	if p, ok := ctx.Get(consts.PARAMETERMAP); ok {
		paramMap = p.(map[string]any)
	}
	if p, ok := ctx.Get(consts.RESULTMAP); ok {
		resultMap = p.(map[string]any)
	}
	for _, predicate := range predicates {
		b = compare(predicate.Key, predicate.Value, predicate.Operator, paramMap, resultMap)
		if (!b && predicateType > 0) || (b && predicateType < 1) {
			return b, nil
		}
	}
	return true, nil
}

func compare(k, v, op string, paramMap, resultMap map[string]any) bool {
	switch op {
	case eq:
		return false
	}
	return true
}
