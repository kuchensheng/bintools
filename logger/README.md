# bintools logger

bintools logger's API 提供了简单的使用接口，支持将日志内容写入控制台、文件等，暂不支持颜色输出

# 特性
1. 低空间占用
2. 支持按日志级别输出到不同的日志文件，并且是按需创建文件 
3. 支持定时分割文件 
4. 支持设置单个文件的最大值 
5. 协程安全

# 安装
```go
go get -u github.com/kuchensheng/bintools/logger
```
# 快速开始
## 简单示例如下：

```go
package main

import "github.com/kuchensheng/bintools/logger"

func main() {
	log := logger.GlobalLogger
	log.Info("你好，%s","库陈胜")
	//Output: 2023-02-17 21:50:14,435 DESKTOP-QO73KDA DESKTOP-QO73KDA  [INFO]  F:/Go_workspace /bintools/logger/options_test.go:57  你好,库陈胜
}
```

## 关联上下文,形成上下链路跟踪
```go
package main

import "context"
import "github.com/kuchensheng/bintools/logger"

func main() {
	log := logger.GlobalLogger
	ctx := context.WithValue(context.TODO(), "TRACEID", "我是traceId")
	log = log.WithContext(ctx)
	log.Info("侬好")
	//Output: 2023-02-17 21:53:04,859 DESKTOP-QO73KDA DESKTOP-QO73KDA 我是traceId [INFO]  F:/Go_workspace/go1.18beta2/src/testing/testing.go:1440  侬好
}
```

## 支持以下日志级别
+ panic (logger.PanicLevel)
+ fatal (logger.Fatalf)
+ error (logger.ErrorLevel)
+ warn (logger.WarnLevel)
+ info (logger.InfoLevel)
+ debug (logger.DebugLevel)
+ trace (logger.TraceLevel)

## 设置全局日志级别
```go
package main

import "github.com/kuchensheng/bintools/logger"

func main() {
	log := logger.GlobalLogger
	log.Level(logger.ErrorLevel)
	log.Level(logger.ErrorLevel)
	log.Info("我不会输出")
	log.Error("我是error，我输出")
	//Output: 2023-02-17 21:59:57,238 DESKTOP-QO73KDA DESKTOP-QO73KDA  [ERROR]  F:/Go_workspace/go1.18beta2/src/testing/testing.go:1440  我是error，我输出
}
```