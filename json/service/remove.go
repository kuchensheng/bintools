package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/executor/js"
	"github.com/kuchensheng/bintools/json/executor/util"
	"github.com/kuchensheng/bintools/json/lib"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func Remove(tenantId, path, method, json string) (int, error) {
	pk := strings.ReplaceAll(path, consts.GlobalPrefix, "")
	pk = strings.ReplaceAll(pk, "/", "")
	pk = pk + strings.ToLower(method)
	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "example", tenantId)
	goFile := filepath.Join(dir, pk+".go")
	if err := os.Remove(goFile); err != nil {
		log.Error().Msgf("无法删除go文件[%s],%v", goFile, err)
		return http.StatusInternalServerError, consts.NewBusinessException(1080500, "无法删除Go脚本文件,文件不存在")
	}
	//移除执行引擎
	lib.RemoveProgramMap(pk)
	if json != "" {
		jsonFile := filepath.Join(dir, json)
		if err := os.Remove(jsonFile); err != nil {
			log.Error().Msgf("无法删除json文件[%s],%v", jsonFile, err)
			return http.StatusInternalServerError, consts.NewBusinessException(1080500, "无法删除Json文件,"+err.Error())
		}
	}
	return http.StatusOK, consts.Ok()
}

func Runner(ctx *gin.Context) (int, error) {
	var content = &model.Script{}
	if data, err := ctx.GetRawData(); err != nil {
		log.Warn().Msgf("无法读取请求体内容,%v", err)
		return http.StatusInternalServerError, consts.NewBusinessException(1080500, "无法读取请求内容:"+err.Error())
	} else if err = json.Unmarshal(data, content); err != nil {
		log.Warn().Msgf("无法解析请求体内容,%v", err)
		return http.StatusBadRequest, consts.NewBusinessException(1080500, "无法解析请求内容:"+err.Error())
	} else {
		switch strings.ToUpper(content.Language) {
		case "GO":
			script := lib.ScriptEngineFunc()
			if res, err := script.Eval(content.Script); err != nil {
				log.Error().Msgf("无法编译成go文件:%v", err)
				return http.StatusBadRequest, consts.NewBusinessException(1080500, "无法编译成go文件："+err.Error())
			} else {
				if content.Method != "" {
					res, err = script.Eval(content.Method)
				}
				if err != nil {
					log.Error().Msgf("无法执行go脚本:%v", err)
					return http.StatusBadRequest, consts.NewBusinessException(1080500, "无法执行go脚本："+err.Error())
				}
				var value any
				if res.Type() == nil {
					value = res.String()
				} else if res.Type() != reflect.TypeOf(func(...any) any { return nil }) {
					return http.StatusBadRequest, consts.NewBusinessException(1080500, "方法["+content.Method+"]类型不对，必须是func(...any) any")
				} else {
					function := res.Interface().(func(...any) any)
					value = function(content.Args)
				}
				return http.StatusOK, consts.OkWithData(value)
			}
		case "JAVASCRIPT":
			step := model.ApixStep{
				Script: model.ApixScript{
					Script:   content.Script,
					Language: "javascript",
				},
				GraphId: "default",
			}
			if err = js.ExecuteJavaScript(ctx, step); err != nil {
				return http.StatusBadRequest, consts.NewBusinessException(1080500, "无法执行script脚本："+err.Error())
			}
			return http.StatusOK, consts.OkWithData(util.GetResultValue(ctx, "$default.$resp.export"))
		default:
			return http.StatusOK, consts.Ok()
		}
	}
	return http.StatusOK, consts.Ok()
}
