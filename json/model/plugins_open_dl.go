//go:build (linux && cgo) || (darwin && cgo) || (freebsd && cgo)

package model

import (
	"github.com/gin-gonic/gin"
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
