package predicate

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/executor/util"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/kuchensheng/bintools/tracer/trace"
	"github.com/rs/zerolog/log"
	"strconv"
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
	tracer, _ := ctx.Get(consts.TRACER)
	clientTracer := tracer.(*trace.ServerTracer).NewClientWithHeader(&ctx.Request.Header)
	clientTracer.TraceName = "执行判断逻辑节点"
	log.Info().Msgf("执行逻辑判断，type= %d", predicateType)
	var b bool

	for _, predicate := range predicates {
		left := util.GetContextValue(ctx, predicate.Key)
		right := util.GetContextValue(ctx, predicate.Value)
		b = compare(left, right, predicate.Operator)
		if (!b && predicateType > 0) || (b && predicateType < 1) {
			return b, nil
		}
	}
	clientTracer.EndTrace(trace.OK, "判断节点执行结果:"+strconv.FormatBool(b))
	return b, nil
}

func compare(k, v any, op string) bool {
	switch op {
	case eq:
		return false
	}
	return true
}
