package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/model"
	"strings"
)

func main() {
	//启动http服务，承接请求
	relativePath := flag.String("context_path", "/api/app/orc", "请求路径前缀")

	router := gin.Default()
	router.POST("/api/app/orc-server/build", func(context *gin.Context) {
		//新增构建
		data, _ := context.GetRawData()
		if _, err := BuildJson(data); err != nil {
			context.JSON(400, model.NewBusinessException(1080500, err.Error()))
		}
	})
	router.Any(*relativePath, func(context *gin.Context) {
		//todo 执行插件
		model.ExecutePlugin(getPluginKey(context), context)
	})
}

func getPluginKey(ctx *gin.Context) string {
	path := ctx.Request.URL.Path
	method := ctx.Request.Method
	method = strings.ToLower(method)
	version := ctx.DefaultQuery("version", "")
	code := ctx.DefaultQuery("code", "")
	if code != "" {
		return code
	}
	key := strings.Join([]string{path, method, version}, "_")
	key = strings.ReplaceAll(key, "/", "")
	if strings.HasPrefix(key, "_") {
		key = strings.ReplaceAll(key, "_", "bintools")
	}
	return key
}
