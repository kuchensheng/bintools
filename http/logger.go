package http

import (
	"github.com/google/uuid"
)

const (
	LoggerKey        = "logger"
	T_HEADER_TRACEID = "t-trace-id"
)

var LoggerMiddleWare = func(ctx *Context) {
	logger := ctx.Logger()
	l := logger.WithContext(ctx)
	traceId := ctx.GetHeader(T_HEADER_TRACEID)
	if traceId == "" {
		traceId = uuid.NewString()
		ctx.SetHeader(T_HEADER_TRACEID, traceId)
	}
	l = l.TraceId(traceId)
	ctx.Set(LoggerKey, l)
	ctx.Next()
}
