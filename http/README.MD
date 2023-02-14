# Isc Web Framework

Isc 是一个用Go（Golang）编写的HTTP web框架，它类似于gin，但比gin更丰富。在IscWeb框架中，使用前缀树进行router分发，分发效率达到 23905 ns/op,可以提供极高的性能。

## Isc Web框架特性
+ 快速：性能极好，单测结果达到 13905 ns/op
+ 简单：使用简单，API友好
+ 容易扩展：支持MiddleWare自由扩展
+ 结果浅封装：对返回结果统一封装，使之具备统一格式
+ 与日志集成，集成了Logger模块，可快速执行日志输出
+ Crash处理，Isc可以catch一个发生在HTTP请求中的panic，并recover，来保证服务使用可用，并将panic信息输出，帮助开发者快速定位原因
+ JSON验证，解析并验证请求的JSON

# 前置条件
+ Go 1.18及以上版本

# 安装
1. 下载并安装Isc
```shell
go get -u github.com/kuchensheng/bintools/http
```
2. 将Isc引入到代码中
```go
import "github.com/kuchensheng/bintools/http"
```

3. 将Logger模块引入代码中（可选）
```go
import "github.com/kuchensheng/bintools/logger"
```

# 运行示例
将以下代码写入 `example.go` 中
```go
package main

import "fmt"
import "github.com/kuchensheng/bintools/http"

func main() {
    e := http.Default()
	e.Get("/api/test/string", func(ctx *http.Context) {
        ctx.JSONoK("成功")
	})
	e.Post("/api/test/struct", func(ctx *http.Context) {
        ctx.JSONoK(http.Result{0,"成功","我是结果"})
	})
	//isc 也支持参数设置
	e.PostWithParam("/api/test/param", func(params ...http.HandlerParam) (any, error) {
        return fmt.Sprintf("%+v", params), nil
	},http.NewQuery("name",false),http.BodyParam{struct {
        Class string `json:"class"`
	}{},false})
	e.Run(8080)
}
```
然后，执行 `go run example.go`，日志提示如下：表示启动成功
```text
2023-02-14 14:13:13,868 DESKTOP-1HAVLBR DESKTOP-1HAVLBR  [INFO]  /home/worksapace/go/bintools/http/engine_test.go:235  服务启动完成，使用端口号:8080
```
在浏览器中访问 http://localhost:8080/api/test/string 即可

# 基准测试
Isc Web framework 使用了自定义的router，基准测试结果如下
  ref
| Benchmark_name       | 个数      | 操作速度        |
|----------------------|---------|-------------|
| BenchmarkEngine_Get  | 	43550	 | 27364 ns/op |
| BenchmarkEngine_Post | 	40543	 | 29670 ns/op |
| BenchmarkAero_Any    | 	57632	 | 20648 ns/op |

# 中间件扩展
只要类型是`func(ctx *Context)` 均可被拓展加入
```go
var AuthMiddleWare = func(ctx *Context) {
	//执行前进行鉴权
	auth = false
	if !auth {
	    ctx.Abort()
		return
    }
	ctx.Next()
}

e := http.Default()
e.Use(LoggerMiddleWare)
e.Run(8080)
```
