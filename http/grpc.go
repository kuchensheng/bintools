package http

import (
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

var grpcServer *grpc.Server
var GrpcContext = func(ctx *Context) {
	if ctx.Request.ProtoMajor == 2 && strings.HasPrefix(ctx.GetHeader("Content-Type"), "application/grpc") {
		ctx.Status(http.StatusOK)
		if grpcServer == nil {
			grpcServer = grpc.NewServer()
		}
		grpcServer.ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Abort()
		return
	}
	//当做普通api
	ctx.Next()
}
