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
	handlers HandlersChain
	pool     sync.Pool
	routes   *Trie
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
	//todo 处理路由规则
	p := pattern
	if !strings.HasPrefix(p, SEP) {
		p = SEP + p
	}
	p = http.MethodGet + p
	e.routes.Insert(p, Route{
		Method:  http.MethodGet,
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
	c := e.pool.Get().(*Context)
	c.Request = req
	c.Writer = w
	//处理该上下文
	method := strings.ToUpper(req.Method)
	uri := req.URL.Path
	word := method + uri
	search := e.routes.Search(word)
	if search == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	c.Route = search
	req.ParseForm()
	c.handlers = search.Handler
	defer c.Recovery()
	handle(c.handlers, c)
	defer e.pool.Put(c)
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
