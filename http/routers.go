package http

type IRoute interface {
	Get(pattern string, handlers ...HandlerFunc)
	Post(pattern string, handlers ...HandlerFunc)
	Put(pattern string, handlers ...HandlerFunc)
	Delete(pattern string, handlers ...HandlerFunc)
	Options(pattern string, handlers ...HandlerFunc)
	Head(pattern string, handlers ...HandlerFunc)
	Any(pattern string, handlers ...HandlerFunc)
}
