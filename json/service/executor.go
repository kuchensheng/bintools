package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var programMap = make(map[string]*interp.Interpreter)
var GoPath = "D:\\worksapace\\go"

var scriptEngineFunc = func() *interp.Interpreter {
	i := interp.New(interp.Options{GoPath: GoPath})
	i.Use(stdlib.Symbols)
	i.Use(Symbols)
	return i
}

func Compile(path string) error {
	key := "default"
	if f, e := os.Stat(path); e != nil {
		log.Error().Msgf("无法打开待编译文件：%v", e)
		return e
	} else {
		key = f.Name()
		key = strings.ReplaceAll(key, ".go_", "")
	}

	log.Info().Msgf("编译文件:%s,key=%s", path, key)
	var scriptEngine = scriptEngineFunc()

	if _, err := scriptEngine.EvalPath(path); err != nil {
		log.Error().Msgf("Go文件无法被编译，%v", err)
		return err
	} else {
		programMap[key] = scriptEngine
	}
	log.Info().Msgf("文件编译完成")
	return nil
}

func Execute(context *gin.Context) (any, error) {
	//执行go脚本
	pk := getPackage(context)
	scriptPath := readGoScript(context, pk)
	var scriptEngine *interp.Interpreter
	if p, ok := programMap[pk]; ok {
		scriptEngine = p
	} else {
		log.Info().Msgf("执行了未编译的脚本,这需要花点时间,pk = %s", pk)
		ch := make(chan error, 1)
		go func() {
			ch <- Compile(scriptPath)
		}()
		select {
		case e := <-ch:
			if e != nil {
				return nil, errors.New("脚本解析失败")
			} else if p, ok = programMap[pk]; ok {
				scriptEngine = p
			}
		case <-time.After(1 * time.Minute):
			log.Warn().Msgf("编译比较耗时，不建议等待")
			return nil, errors.New("正在执行脚本初始化，请稍后再试")
		}
		log.Info().Msgf("脚本编译完成")
	}
	if v, e := scriptEngine.Eval(fmt.Sprintf("%s.%s%s", pk, "Executor", pk)); e != nil {
		log.Error().Msgf("Go 脚本编译异常,%v", e)
	} else if v.Type() != nil {
		return v.Interface().(func(ctx *gin.Context) (any, error))(context)
	}

	return nil, nil

}

func getPackage(ctx *gin.Context) string {
	uri := strings.ReplaceAll(ctx.Request.URL.Path, "/", "_")
	method := strings.ToLower(ctx.Request.Method)
	version := ctx.GetHeader("version")
	key := strings.Join([]string{uri, method, version}, "_")
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "_api_app_orc_", "")
	key = strings.ReplaceAll(key, "_", "")
	return key
}

func readGoScript(ctx *gin.Context, key string) string {
	wd, _ := os.Getwd()
	tenantId := ctx.GetHeader("isc-tenant-id")
	fp := filepath.Join(wd, "example", tenantId, key+".go_")
	return fp
}
