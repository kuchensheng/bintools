package model

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"plugin"
)

var plugins map[string]*PluginInfo

const (
	PluginDll = iota
	PluginSo
)

type PluginInfo struct {
	plugin.Plugin
	plugin.Symbol
	PluginType int
}

type PluginDefinition struct {
	Name   string `json:"name"`   //插件名称
	Path   string `json:"path"`   //插件路径（绝对路径）
	Key    string `json:"key"`    //插件唯一标识
	Method string `json:"method"` //插件执行入口方法

}

func ExecutePlugin(key string, ctx *gin.Context) {
	GetPluginByKey(key).Execute(ctx)
}

func GetPluginByKey(key string) *PluginInfo {
	return plugins[key]
}

func LoadPlugin(definition PluginDefinition) error {
	if info, err := OpenPlugin(definition); err != nil {
		log.Error().Msgf("无法打开插件%s,%v", definition.Name, err)
		return err
	} else {
		plugins[definition.Key] = info
	}
	//todo 持久化处理
	return nil
}
