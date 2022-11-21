//go:build windows

package model

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
	defer func() {
		if e := recover(); e != nil {
			context.JSON(400, NewBusinessExceptionWithData(1080500, "执行规则执行异常", e))
		}
	}()
	proc := plugin.Symbol.(*syscall.LazyProc)
	proc.Call(uintptr(unsafe.Pointer(context.Request)))
}

func (definition PluginDefinition) Execute(ctx *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Error().Msgf("插件执行异常,%v", x)
		}
	}()
	requestPoint := uintptr(unsafe.Pointer(ctx.Request))
	writer := &ctx.Writer
	writerPoint := uintptr(unsafe.Pointer(writer))
	library, err := syscall.LoadLibrary(definition.Path)
	//defer syscall.FreeLibrary(library)
	if err != nil {
		log.Error().Msgf("无法加载插件，%v", err)
		ctx.JSON(400, NewBusinessExceptionWithData(1080500, "无法加载插件", err))
		return
	}
	proc, err := syscall.GetProcAddress(library, definition.Method)
	if err != nil {
		log.Error().Msgf("无法打开插件，%v", err)
		ctx.JSON(400, NewBusinessExceptionWithData(1080500, "无法打开插件", err))
		return
	}
	if r1, r2, err1 := syscall.SyscallN(proc, requestPoint, writerPoint); err1 != syscall.Errno(0) {
		println("执行结束", r1, r2, err1)
	}
}
