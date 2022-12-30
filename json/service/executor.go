package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	log2 "github.com/kuchensheng/bintools/json/executor/log"
	"github.com/kuchensheng/bintools/json/lib"
	"github.com/rs/zerolog/log"
	"github.com/traefik/yaegi/interp"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

func Compile(path string) error {
	key := "default"
	if f, e := os.Stat(path); e != nil {
		log.Error().Msgf("无法打开待编译文件：%v", e)
		return e
	} else {
		key = f.Name()
		key = strings.ReplaceAll(key, ".go", "")
	}

	log.Info().Msgf("编译文件:%s,key=%s", path, key)
	var scriptEngine = lib.ScriptEngineFunc()

	if _, err := scriptEngine.EvalPath(path); err != nil {
		log.Error().Msgf("Go文件无法被编译，%v", err)
		return err
	} else {
		lib.PutProgramMap(key, scriptEngine)
	}
	log.Info().Msgf("文件编译完成")
	return nil
}

func Execute(ctx *gin.Context) {
	ch := make(chan error, 1)
	var result any
	go func(channel chan error, ctx *gin.Context) {
		//获取请求体
		//r, e := bweditpost.Executorbweditpost(ctx)
		r, e := execute(ctx)
		channel <- e
		result = r
	}(ch, ctx)
	select {
	case err := <-ch:
		if err != nil {
			ctx.Set(consts.ErrKey, err)
			if reflect.TypeOf(err) == reflect.TypeOf(consts.NewException("", "", "")) {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				ctx.JSON(http.StatusBadRequest, consts.NewBusinessException(1080400, err.Error()))
			}
			return
		}
	case <-time.After(30 * time.Second):
		ctx.JSON(400, consts.NewBusinessException(1080500, "请求超时请检查"))
		return
	}
	ctx.JSON(http.StatusOK, result)
	return
}

func execute(context *gin.Context) (any, error) {
	defer func() {
		if x := recover(); x != nil {
			log.Error().Msgf("请求执行异常，panic:%v", x)
			fmt.Sprintf("%s\n", debug.Stack())
		}
	}()
	//执行go脚本
	pk := log2.GetPackage(context)
	scriptPath := readGoScript(context, pk)
	var scriptEngine *interp.Interpreter
	if p, ok := lib.GetProgramMap(pk); ok {
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
			} else if p, ok = lib.GetProgramMap(pk); ok {
				scriptEngine = p
			}
		case <-time.After(1 * time.Minute):
			log.Warn().Msgf("编译比较耗时，不建议等待")
			return nil, errors.New("正在执行脚本初始化，请稍后再试")
		}
		scriptEngine, _ = lib.GetProgramMap(pk)
		log.Info().Msgf("脚本编译完成")
	}
	if v, e := scriptEngine.Eval(fmt.Sprintf("%s.%s%s", pk, "Executor", pk)); e != nil {
		log.Error().Msgf("Go 脚本编译异常,%v", e)
	} else if v.Type() != nil {
		return v.Interface().(func(ctx *gin.Context) (any, error))(context)
	}

	return nil, nil

}

func readGoScript(ctx *gin.Context, key string) string {
	wd, _ := os.Getwd()
	tenantId := ctx.GetHeader("isc-tenant-id")
	fp := filepath.Join(wd, "example", tenantId, key+".go")
	return fp
}
