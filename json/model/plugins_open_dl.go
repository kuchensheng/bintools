//go:build (linux && cgo) || (darwin && cgo) || (freebsd && cgo)

package model

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/example"
	"net/http"
	"plugin"
)

func OpenPlugin(definition PluginDefinition) (*PluginInfo, error) {
	p, err := plugin.Open(definition.Path)
	if err != nil {
		return nil, err
	}
	symbol, err := p.Lookup(definition.Method)
	if err != nil {
		return nil, err
	}
	pp := &PluginInfo{
		Plugin:     *p,
		Symbol:     symbol,
		PluginType: PluginSo,
	}
	return pp, nil
}

func (plugin *PluginInfo) Execute(context *gin.Context) {
	symbol := plugin.Symbol
	symbol.(func(r *http.Request, w *http.Response))(context.Request, context.Request.Response)
}

func (definition PluginDefinition) Execute(context *gin.Context) {
	p, err := plugin.Open(definition.Path)
	if err != nil {
		context.JSON(400, NewBusinessExceptionWithData(1080500, "无法执行对应规则", err))
		return
	}
	if _, err := p.Lookup(definition.Method); err != nil {
		context.JSON(400, NewBusinessExceptionWithData(1080500, "无法执行对应规则", err))
		return
	} else {
		example.Executor(context.Request, context.Writer)
		//symbol.(func(r *http.Request, w http.ResponseWriter))(context.Request, context.Writer)
	}
}
