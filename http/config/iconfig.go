package config

var (
	ConfigMap map[string]any
	Concat    = "."
)

type Config struct {
}

type IAppConfig interface {

	//GetAttr 获取属性值，key = a.b.c
	GetAttr(key string) any

	//FillAttr 给对象字段填充之
	FillAttr(pointer any, prefix string)
}
