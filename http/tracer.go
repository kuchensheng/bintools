package http

import trace2 "github.com/kuchensheng/bintools/tracer/trace"

var TracerMiddleWare = func(context *Context) {
	//服务端接入tracer
	tracerS := trace2.NewServerTracer(context.Request)
	//将结果写入context即可
	context.Next()
	//结束时
	tracerS.EndServerTracer(trace2.OK, "")
}
