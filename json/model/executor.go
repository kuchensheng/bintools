package model

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/extractlib"
	"github.com/rs/zerolog/log"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var programMap = make(map[string]*interp.Interpreter)

var scriptEngineFunc = func() *interp.Interpreter {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	i.Use(extractlib.Symbols)
	return i
}

func Compile(path, key string) error {
	log.Info().Msgf("编译文件:%s,key=%s", path, key)
	var scriptEngine = scriptEngineFunc()

	if _, err := scriptEngine.EvalPath(path); err != nil {
		log.Error().Msgf("Go文件无法被编译，%v", err)
		return err
	} else {
		programMap[key] = scriptEngine
	}
	return nil
}

type fv func(ctx *gin.Context) (any, error)

var _fv = fv(nil)

func Execute(context *gin.Context) (any, error) {
	//执行go脚本
	pk := getPackage(context)
	pk = strings.ReplaceAll(pk, "_api_app_orc_", "")
	pk = strings.ReplaceAll(pk, "_", "")
	script := readGoScript(context, pk)

	var scriptEngine *interp.Interpreter
	var err error
	if p, ok := programMap[pk]; ok {
		scriptEngine = p
	} else {
		scriptEngine = scriptEngineFunc()
		if _, err = scriptEngine.Eval(script); err != nil {
			log.Error().Msgf("脚本解析异常,%v", err)
			return nil, err
		}
	}

	v, _ := scriptEngine.Eval("bweditpost.Executor")
	fu := v.Interface().(func(ctx *gin.Context) (any, error))
	return fu(context)

}

func getPackage(ctx *gin.Context) string {
	uri := strings.ReplaceAll(ctx.Request.URL.Path, "/", "_")
	method := strings.ToLower(ctx.Request.Method)
	version := ctx.GetHeader("version")
	key := strings.Join([]string{uri, method, version}, "_")
	key = strings.ReplaceAll(key, "/", "_")
	return key
}

func readGoScript(ctx *gin.Context, key string) string {
	wd, _ := os.Getwd()

	fp := filepath.Join(wd, "example", key+".go_")
	if data, err := ioutil.ReadFile(fp); err != nil {
		log.Error().Msgf("文件读取失败,path=[%s],错误信息：%s", fp, err)
		return ""
	} else {
		return string(data)
	}
}
