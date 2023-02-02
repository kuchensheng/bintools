package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type engine struct {
	//全局执行调用链
	handlers HandlersChain
	// 上下文池子，避免重复释放-申请内存
	pool sync.Pool
	//路由规则前缀树
	routes *Trie
}

func Default() *engine {
	return &engine{
		handlers: HandlersChain{},
		pool: sync.Pool{
			New: func() any {
				return &Context{}
			},
		},
		routes: NewTrie(),
	}
}

func (e *engine) Use(middlewares ...HandlerFunc) {
	e.handlers = append(e.handlers, middlewares...)
}

//Get 注册路由规则及执行方法
func (e *engine) Get(pattern string, handlers ...HandlerFunc) {
	e.registerRouter(http.MethodGet, pattern, handlers...)
}

func (e *engine) Post(pattern string, handlers ...HandlerFunc) {
	e.registerRouter(http.MethodPost, pattern, handlers...)
}
func (e *engine) Put(pattern string, handlers ...HandlerFunc) {
	e.registerRouter(http.MethodPut, pattern, handlers...)
}
func (e *engine) Delete(pattern string, handlers ...HandlerFunc) {
	e.registerRouter(http.MethodDelete, pattern, handlers...)
}
func (e *engine) Options(pattern string, handlers ...HandlerFunc) {
	e.registerRouter(http.MethodOptions, pattern, handlers...)
}
func (e *engine) Head(pattern string, handlers ...HandlerFunc) {
	e.registerRouter(http.MethodHead, pattern, handlers...)
}
func (e *engine) Any(pattern string, handlers ...HandlerFunc) {
	e.registerRouter(http.MethodGet, pattern, handlers...)
	e.registerRouter(http.MethodPost, pattern, handlers...)
	e.registerRouter(http.MethodPut, pattern, handlers...)
	e.registerRouter(http.MethodDelete, pattern, handlers...)
	e.registerRouter(http.MethodOptions, pattern, handlers...)
	e.registerRouter(http.MethodHead, pattern, handlers...)
}

func (e *engine) registerRouter(method, pattern string, handlers ...HandlerFunc) {
	p := pattern
	if !strings.HasPrefix(p, SEP) {
		p = SEP + p
	}
	p = method + p
	e.routes.Insert(p, Route{
		Method:  method,
		Path:    pattern,
		Handler: append(e.handlers, handlers...),
	})
}

func (e *engine) Run(port int) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: e,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("unable to start server due to: %v", err)
		}
	}()
	//优雅停机
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGILL, syscall.SIGABRT, syscall.SIGSEGV)
	<-quit
	log.Print("Shutdown Server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error , %v", err)
	}
	select {
	case <-ctx.Done():
	}
	log.Print("Server exiting")
}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//查找目标处理器
	search := func(request *http.Request) *Route {
		method := strings.ToUpper(request.Method)
		uri := req.URL.Path
		word := method + uri
		return e.routes.Search(word)
	}(req)
	if search == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	c := e.pool.Get().(*Context)
	c.reset()
	c.Request = req
	c.Writer = w

	c.Route = search
	req.ParseForm()
	c.handlers = search.Handler
	defer c.Recovery()
	handle(c.handlers, c)
	defer e.pool.Put(c)
}

func (c *Context) reset() {
	c.index = 0
	c.Writer = nil
	c.Request = nil
	c.handlers = nil
	c.Keys = nil
	c.queryCache = nil
	c.formCache = nil
	c.Route = nil
}

func handle(chain HandlersChain, ctx *Context) {
	ctx.handlers[ctx.index](ctx)
	//最后一个必须得执行
	if ctx.index < int8(len(ctx.handlers)-1) {
		handle(chain[ctx.index+1:], ctx)
	}
}

// Last returns the last handler in the chain. ie. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
