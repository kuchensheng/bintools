package http

import (
	"context"
	"fmt"
	"github.com/kuchensheng/bintools/http/config"
	"github.com/kuchensheng/bintools/http/config/yaml"
	"github.com/kuchensheng/bintools/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"reflect"
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
	r := &Route{method, pattern, h, nil, nil}
	e.routes.Insert(p, r)
	logger.GlobalLogger.Info("注册路由规则:Method [%s],Pattern [%s]", method, pattern)
}

func (e *engine) registerRouterWithParams(method, pattern string, handler HandlerParamFunc) {
	p := pattern
	if !strings.HasPrefix(p, SEP) {
		p = SEP + p
	}
	p = method + p
	h := make(HandlersChain, len(e.handlers))
	_ = copy(h, e.handlers)
	r := &Route{method, pattern, h, handler, func() []string {
		var res []string
		reflect.TypeOf(handler)
		return res
	}()}
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
	if search.ParamsHandler != nil {
		if handler, err := search.ParamsHandler(parseParams(search, c)); err != nil {
			c.JSON(http.StatusBadRequest, e)
		} else {
			c.JSON(http.StatusOK, handler)
		}
	}
	defer e.pool.Put(c)
}

func parseParams(route *Route, context *Context) []HandlerParam {
	var res []HandlerParam
	split := strings.Split(route.Path, "/:")
	if len(split) == 1 {
		split = strings.Split(route.Path, "/{")
	}
	if len(split) == 1 {
		//当前请求无Path参数
	} else {
		split = strings.Split(split[1], "/")
		pathKey := split[0]
		//获取pathKey所在idx
		for idx, s := range strings.Split(route.Path, "/") {
			if s == pathKey {
				pathVal := strings.Split(context.Request.RequestURI, "/")[idx]
				res = append(res, PathParam{name: pathKey, value: pathVal})
				break
			}
		}
	}

	//查询请求头
	for s, vals := range context.Request.Header {
		res = append(res, HeaderParam{name: s, value: vals[0]})
	}
	//查询query参数
	for key, values := range context.Request.URL.Query() {
		res = append(res, QueryParam{name: key, value: values[0]})
	}
	//查询Body参数
	if data, err := context.GetRawData(); err == nil {
		res = append(res, BodyParam{Body: data})
	}
	//查询表单参数
	if context.GetHeader("Content-Type") == "application/x-www-form-urlencoded" {
		if context.Request.PostForm != nil {
			res = append(res, FormParam{Form: context.Request.PostForm})
		} else {
			res = append(res, FormParam{Form: context.Request.Form})
		}
	}
	//查询文件表单
	if context.GetHeader("Content-Type") == "multipart/form-data" {
		if context.Request.MultipartForm != nil {
			res = append(res, MultiFormParam{Form: context.Request.MultipartForm})
		}
	}
	return res
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
