//go:build windows

package model

import (
	"github.com/gin-gonic/gin"
	"syscall"
	"unsafe"
)

func OpenPlugin(definition PluginDefinition) (*PluginInfo, error) {
	dll := syscall.NewLazyDLL(definition.Path)
	proc := dll.NewProc(definition.Method)
	return &PluginInfo{
		Symbol:     proc,
		PluginType: PluginDll,
	}, nil
}

func (plugin *PluginInfo) Execute(context *gin.Context) {
	proc := plugin.Symbol.(*syscall.LazyProc)
	_, _, _ := proc.Call(uintptr(unsafe.Pointer(context.Request)), uintptr(unsafe.Pointer(context.Request.Response)))
}
