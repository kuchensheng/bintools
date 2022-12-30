package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/tracer/trace"
)

func TracerFilter() gin.HandlerFunc {
	return func(context *gin.Context) {
		tracer := trace.NewServerTracer(context.Request)
		context.Set(consts.TRACER, tracer)
		context.Next()
		if v, ok := context.Get(consts.ErrKey); ok {
			tracer.EndServerTracer(trace.WARNING, v.(error).Error())
		} else {
			tracer.EndServerTracer(trace.OK, "")
		}
	}
}
