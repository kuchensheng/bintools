package predicate

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	log2 "github.com/kuchensheng/bintools/json/executor/log"
	"github.com/kuchensheng/bintools/json/executor/util"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/kuchensheng/bintools/tracer/trace"
	"github.com/rs/zerolog/log"
	"regexp"
	"strconv"
	"strings"
)

const (
	eq     string = "=="
	ne     string = "!="
	ge     string = ">="
	gt     string = ">"
	le     string = "<="
	lt     string = "<"
	inc    string = "inc"
	_inc   string = "!inc"
	_nil   string = "nil"
	__nil  string = "!nil"
	_true  string = "true"
	_false string = "false"
)

//ExecPredicates predicateType = 1 所有条件都为真，否则任一条件为真
func ExecPredicates(ctx *gin.Context, step model.ApixStep) (bool, error) {
	tracer, _ := ctx.Get(consts.TRACER)
	clientTracer := tracer.(*trace.ServerTracer).NewClientWithHeader(&ctx.Request.Header)
	pk := log2.GetPackage(ctx)
	ls := log2.LogStruct{PK: pk, TraceId: clientTracer.TracId}
	clientTracer.TraceName = "执行判断逻辑节点"
	ls.Info(clientTracer.TraceName + "...")
	predicateType := step.PredicateType
	log.Info().Msgf("执行逻辑判断节点：%s，type= %d", step.GraphId, predicateType)
	var b bool
	predicates := step.Predicate
	for _, predicate := range predicates {
		left := util.GetContextValue(ctx, predicate.Key)
		right := util.GetContextValue(ctx, predicate.Value)
		ls.Info("逻辑判断比较,%s %s %s", fmt.Sprintf("%v", left), predicate.Operator, fmt.Sprintf("%v", right))
		b = compare(left, right, predicate.Operator)
		if !b && predicateType > 0 {
			break
		} else if b && predicateType < 1 {
			break
		}
	}
	ls.Info("逻辑判断执行完毕,执行结果：%t", b)
	log.Info().Msgf("逻辑判断执行完毕,执行结果：%t", b)
	clientTracer.EndTrace(trace.OK, "判断节点执行结果:"+strconv.FormatBool(b))
	return b, nil
}

func compare(k, v any, op string) bool {
	log.Info().Msgf("kv比较,k=%s,op = %s,v=%s", k, op, v)
	if k == nil {
		log.Warn().Msgf("k不能为空,返回false")
		return false
	}
	defer func() {
		if x := recover(); x != nil {
			log.Error().Msgf("逻辑判断出错，%v", x.(error))
		}
	}()
	switch op {
	case eq:
		return fmt.Sprintf("%v", k) == fmt.Sprintf("%v", v)
	case inc:
		return strings.Contains(k.(string), v.(string))
	case ne:
		return k != v
	case ge:
		return k.(int64) >= v.(int64)
	case gt:
		return k.(int64) > v.(int64)
	case le:
		return k.(int64) <= v.(int64)
	case lt:
		return k.(int64) < v.(int64)
	case _inc:
		return !strings.Contains(k.(string), v.(string))
	case _nil:
		return v == nil
	case __nil:
		return v != nil
	case _true:
		fallthrough
	case _false:
		res, _ := regexp.Match(k.(string), v.([]byte))
		return res
	default:
		return true
	}

	return true
}
