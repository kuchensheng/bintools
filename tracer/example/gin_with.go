package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/tracer/trace"
)

func main() {

	router := gin.New()
	router.Use(func(context *gin.Context) {
		//服务端接入tracer
		tracerS := trace.NewServerTracer(context.Request)
		//将结果写入context即可
		context.Next()
		//结束时
		tracerS.EndServerTracer(trace.TraceStatusEnum(context.GetInt(trace.T_RESULT_CODE)), context.GetString(trace.T_RESULT_MSG))
	})
}
