# 简介
链路跟踪isc-gobase-tracer为分布式应用提供了完整的调用链路还原、调用请求量统计、链路拓扑、应用依赖分析等工具，可以帮助开发者快速分析和诊断分布式应用架构下的性能瓶颈，提供微服务时代下的开发诊断效率
# 主要功能
+ 分布式调用链查询和诊断：追踪分布式架构中的所有微服务用户请求，并将它们汇总成分布式调用链
+ 分布式拓扑动态发现：用户的所有分布式微服务应用和相关产品可以通过链路追踪收集到分布式调用信息
+ 丰富的下游对接场景：收集的链路可直接用于日志分析，且可对接到系统管理-运维中心等下游分析平台。
# 快速开始
## Install
```bash
go get github.com/isyscore/isc-gobase/tracer
```
## Example
```go
//初始化配置信息
var Conf = &ServiceConf{
	//当前服务名
    ServiceName: "default",
	//链路跟踪保存策略，默认Loki保存
    Using:       "loki",
	//Loki保存策略的配置信息
    Loki: lokiConf{
		//Loki地址
        Host:        "http://loki-service:3100",
		//批量提交的最大值，默认64条
        MaxBatch:    512,
		//提交前最大的等待时间，单位秒，默认1秒
        MaxWaitTime: 1,
        },
}
//create a server tracer
"github.com/isyscore/isc-gobase/tracer/conf"
"github.com/isyscore/isc-gobase/tracer/push"
func testReq(req *http.Request)  {
	//开启客户端跟踪
    serverTracer := NewServerTracer(req)
    println("服务端其他业务请求")
    
    for i := 0; i < 3; i++ {
        println("作为客户端，向其他服务发起请求")
		req1 := &http.Request{}
        clientTracer := serverTracer.NewClientTracer(req1)
		println("req1请求处理以及其他业务处理")
		//结束当前客户端请求跟踪
        clientTracer.EndTrace(OK, "i am danger")
    }
	//结束服务端跟踪
    serverTracer.EndTrace(OK, "i am not in danger")
}
```
# 数据如何上报？
