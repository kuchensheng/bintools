package main

import (
	"github.com/kuchensheng/bintools/http"
	"github.com/kuchensheng/bintools/trace/trace"
	http2 "net/http"
)

func main() {
	e := http.Default()
	e.Use(func(context *http.Context) {
		//服务端接入tracer
		tracerS := trace.NewServerTracer(context.Request)
		//将结果写入context即可
		context.Next()
		//结束时
		tracerS.EndServerTracer(trace.TraceStatusEnum(context.GetInt(trace.T_RESULT_CODE)), context.GetString(trace.T_RESULT_MSG))
	})
	e.Get("/api/test/", func(ctx *http.Context) {
		server := trace.NewServerTracer(ctx.Request)
		req, _ := http2.NewRequest(http2.MethodGet, "http://www.baidu.com", nil)
		server.Call(req)
		defer server.EndTraceOk()
		ctx.JSONoK("你好")
	})
	e.Run(8080)
}
