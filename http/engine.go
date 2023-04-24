package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kuchensheng/bintools/http/config"
	"github.com/kuchensheng/bintools/http/config/yaml"
	"github.com/kuchensheng/bintools/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
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
	routes *trie

	RpcServer bool

	pprof bool

	appConfig config.IAppConfig
}

func Default() *engine {
	e := &engine{
		//handlers: HandlersChain{LoggerMiddleWare, GrpcContext},
		pool: sync.Pool{
			New: func() any {
				return &Context{}
			},
		},
		routes: NewTrie(),
		pprof:  false,
	}
	e.Use(LoggerMiddleWare, GrpcContext)
	e.appConfig = yaml.InitYamlConfig()
	println(config.ConfigMap)
	return e
}

func (e *engine) Pprof(open bool) {
	e.pprof = open
}

func (e *engine) Use(middlewares ...HandlerFunc) {
	e.handlers = append(e.handlers, middlewares...)
}

func handlerParam2Fun(paramFunc HandlerParamFunc, params ...HandlerParam) func(ctx *Context) {
	return func(ctx *Context) {
		log := ctx.Logger()
		var paramValues []HandlerParam
		for _, param := range params {
			if q, ok := param.(QueryParam); ok {
				if v, find := ctx.GetQuery(q.Name()); !find && q.Required() {
					log.Info("缺少请求参数：%s", q.Name())
					e := BADREQUEST
					e.Data = "缺少请求参数:" + q.Name()
					ctx.JSON(http.StatusBadRequest, e)
					ctx.Next()
					return
				} else {
					q.value = v
					paramValues = append(paramValues, q)
				}
			} else if b, ok := param.(BodyParam); ok {
				data, _ := ctx.GetRawDataNoClose()
				body := b.Body
				if len(data) == 0 {
					if b.Required() || ctx.Request.Method == http.MethodPost || ctx.Request.Method == http.MethodPut {
						e := BADREQUEST
						e.Data = "请求体不能为空"
						ctx.JSON(http.StatusBadRequest, e)
						ctx.Abort()
						return
					}
				} else if err := json.Unmarshal(data, &body); err != nil {
					ctx.JSON(http.StatusBadRequest, err)
					ctx.Abort()
					return
				} else {
					b.Body = body
					paramValues = append(paramValues, b)
				}
			}
		}
		if result, err := paramFunc(paramValues...); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSONoK(result)
		}
	}
}

func (e *engine) GetWithParam(pattern string, paramFunc HandlerParamFunc, params ...HandlerParam) {
	e.Get(pattern, handlerParam2Fun(paramFunc, params...))
}

func (e *engine) PostWithParam(pattern string, paramFunc HandlerParamFunc, params ...HandlerParam) {
	e.Post(pattern, handlerParam2Fun(paramFunc, params...))
}

func (e *engine) PutWithParam(pattern string, paramFunc HandlerParamFunc, params ...HandlerParam) {
	e.Put(pattern, handlerParam2Fun(paramFunc, params...))
}

func (e *engine) DeleteWithParam(pattern string, paramFunc HandlerParamFunc, params ...HandlerParam) {
	e.Delete(pattern, handlerParam2Fun(paramFunc, params...))
}

func (e *engine) OptionsWithParam(pattern string, paramFunc HandlerParamFunc, params ...HandlerParam) {
	e.Options(pattern, handlerParam2Fun(paramFunc, params...))
}

func (e *engine) HeadWithParam(pattern string, paramFunc HandlerParamFunc, params ...HandlerParam) {
	e.Head(pattern, handlerParam2Fun(paramFunc, params...))
}

func (e *engine) AnyWithParam(pattern string, paramFunc HandlerParamFunc, params ...HandlerParam) {
	e.Any(pattern, handlerParam2Fun(paramFunc, params...))
}

// Get 注册路由规则及执行方法
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

func (e *engine) Static(relativePath, root string) {
	e.StaticFS(relativePath, Dir(root, true))
}

func (e *engine) StaticFS(relativePath string, fs http.FileSystem) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	//创建handler
	handler := func(r string, f http.FileSystem) HandlerFunc {
		fileServer := http.StripPrefix("/", http.FileServer(fs))
		return func(ctx *Context) {
			if _, noListing := fs.(*onlyFilesFS); noListing {
				ctx.Writer.WriteHeader(http.StatusNotFound)
			}
			filepath, _ := ctx.GetPath("filepath")
			if file, err := fs.Open(filepath); err != nil {
				ctx.Writer.WriteHeader(http.StatusNotFound)
				ctx.Abort()
				return
			} else {
				file.Close()
				fileServer.ServeHTTP(ctx.Writer, ctx.Request)
			}

		}
	}(relativePath, fs)

	urlPattern := path.Join(relativePath, "/*filepath")

	e.registerRouter(http.MethodGet, urlPattern, handler)
	e.registerRouter(http.MethodHead, urlPattern, handler)
}

func (e *engine) registerRouter(method, pattern string, handlers ...HandlerFunc) {
	p := pattern
	if !strings.HasPrefix(p, SEP) {
		p = SEP + p
	}
	p = method + p
	h := make(HandlersChain, len(e.handlers))
	_ = copy(h, e.handlers)
	h = append(h, handlers...)
	r := &Route{method, pattern, h}
	e.routes.Insert(p, r)
	logger.GlobalLogger.Info("注册路由规则:Method [%s],Pattern [%s]", method, pattern)
}

func (e *engine) Run() {
	e.RunWithPort(e.appConfig.GetAttr("server.port").(int))
}

func (e *engine) RunWithPort(port int) {
	if e.pprof {
		e.pprofRouteRegister()
	}
	e.RunRpc(port + 1)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: e,
	}
	l := logger.GlobalLogger

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf(fmt.Sprintf("unable to start server due to: %v", err))
		}
	}()
	l.Info("服务启动完成，使用端口号:%d", port)
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
	if ctx.index < int8(len(ctx.handlers)) {
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
