package model

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
)

var plugins = make(map[string]*PluginInfo)
var pluginDefinitions = make(map[string]PluginDefinition)

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

func GetAllPlugins() map[string]*PluginInfo {
	return plugins
}

func ExecutePlugin(key string, ctx *gin.Context) {
	if pluginDefinition, ok := pluginDefinitions[key]; ok {
		//打开plugin
		pluginDefinition.Execute(ctx)
	} else {
		//从文件中直接读取
		wd, _ := os.Getwd()
		suffix := ".so"
		if runtime.GOOS == "windows" {
			suffix = "_windows.dll"
		}
		pluginPath := filepath.Join(wd, "plugins", key+suffix)
		if _, err := os.Stat(pluginPath); err != nil {
			ctx.JSON(400, NewBusinessException(1080404, "找不到对应的执行规则"))
			return
		}
		//todo 啥也不做，第二次会panic
		d := PluginDefinition{
			Path:   pluginPath,
			Key:    key,
			Method: "Executor",
		}
		pluginDefinitions[key] = d

		d.Execute(ctx)
	}
}

func GetPluginByKey(key string) *PluginInfo {
	return plugins[key]
}

func LoadPlugin(definition PluginDefinition) error {
	pluginDefinitions[definition.Key] = definition
	if info, err := OpenPlugin(definition); err != nil {
		log.Error().Msgf("无法打开插件%s,%v", definition.Name, err)
		return err
	} else {
		plugins[definition.Key] = info
	}
	//todo 持久化处理
	return nil
}
