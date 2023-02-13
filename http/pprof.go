package http

import (
	"net/http"
	"net/http/pprof"
)

const DefaultPrefix = "/debug/pprof"

func (e *engine) pprofRouteRegister() {
	e.Get(getPrefix("/"), pprofHandler(pprof.Index))
	e.Get(getPrefix("/cmdline"), pprofHandler(pprof.Cmdline))
	e.Get(getPrefix("/profile"), pprofHandler(pprof.Profile))
	e.Post(getPrefix("/symbol"), pprofHandler(pprof.Symbol))
	e.Get(getPrefix("/symbol"), pprofHandler(pprof.Symbol))
	e.Get(getPrefix("/trace"), pprofHandler(pprof.Trace))
	e.Get(getPrefix("/allocs"), pprofHandler(pprof.Handler("allocs").ServeHTTP))
	e.Get(getPrefix("/block"), pprofHandler(pprof.Handler("block").ServeHTTP))
	e.Get(getPrefix("/goroutine"), pprofHandler(pprof.Handler("goroutine").ServeHTTP))
	e.Get(getPrefix("/heap"), pprofHandler(pprof.Handler("heap").ServeHTTP))
	e.Get(getPrefix("/mutex"), pprofHandler(pprof.Handler("mutex").ServeHTTP))
	e.Get(getPrefix("/threadcreate"), pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
}

func pprofHandler(h http.HandlerFunc) HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(ctx *Context) {
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func getPrefix(pattern string) string {
	return DefaultPrefix + pattern
}
