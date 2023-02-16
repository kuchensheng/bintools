package trace

type GenerateTraceId interface {
	//GenerateTraceId 生成或获取到唯一traceId值
	GenerateTraceId() string
}
