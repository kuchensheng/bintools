package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io/ioutil"
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
	script := readGoScript(context, pk)
	var scriptEngine *interp.Interpreter
	var err error
	if p, ok := programMap[pk]; ok {
		scriptEngine = p
	} else {
		log.Info().Msgf("执行了未编译的脚本,这需要花点时间,pk = %s", pk)
		ch := make(chan *interp.Interpreter, 1)
		go func() {
			scriptEngine = scriptEngineFunc()
			if _, err = scriptEngine.Eval(script); err != nil {
				log.Error().Msgf("脚本解析异常,%v", err)
				ch <- nil
			} else {
				ch <- scriptEngine
			}
		}()
		select {
		case scriptEngine = <-ch:
			if scriptEngine != nil {
				programMap[pk] = scriptEngine
			} else {
				return nil, errors.New("脚本解析失败")
			}
		case <-time.After(5 * time.Second):
			log.Warn().Msgf("编译比较耗时，不建议等待")
			return nil, errors.New("正在执行脚本初始化，请稍后再试")
		}
		log.Info().Msgf("脚本编译完成")
	}
	v, _ := scriptEngine.Eval(fmt.Sprintf("%s.%s%s", pk, "Executor", pk))
	fu := v.Interface().(func(ctx *gin.Context) (any, error))
	return fu(context)
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
	if data, err := ioutil.ReadFile(fp); err != nil {
		log.Error().Msgf("文件读取失败,path=[%s],错误信息：%s", fp, err)
		return ""
	} else {
		return string(data)
	}
}
