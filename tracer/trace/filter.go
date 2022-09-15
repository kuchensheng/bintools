package trace

import "net/http"

func ServerTraceHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		//开始tracer
		serverTracer := NewServerTracer(r)
		defer func(tracer *ServerTracer) {
			if x := recover(); x != nil {
				tracer.EndTraceError(x.(error))
			} else {
				tracer.EndTraceOk()
			}
		}(serverTracer)
		//这里是业务逻辑
		handler.ServeHTTP(wr, r)
	}
}
