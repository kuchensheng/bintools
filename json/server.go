package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/lib"
	"github.com/kuchensheng/bintools/json/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {
	//启动http服务，承接请求
	relativePath := flag.String("context_path", "/api/app/orc/", "请求路径前缀")
	serverPort := flag.Int("port", 38240, "服务器端口，默认:38240")
	goPath := flag.String("go_path", "", "Go编译环境地址")
	flag.Parse()
	if *goPath != "" {
		lib.GoPath = *goPath
	}
	if *relativePath != consts.GlobalPrefix {
		consts.GlobalPrefix = *relativePath
	}
	log.Logger = log.Logger.Level(zerolog.InfoLevel)
	wd, _ := os.Getwd()
	router := gin.Default()
	router.POST("/api/app/orc-server/build", func(context *gin.Context) {
		data, _ := context.GetRawData()
		tenantId := context.GetHeader(consts.TENANT_ID)
		if _, err := service.BuildJson(data, tenantId); err != nil {
			context.JSON(400, consts.NewBusinessException(1080500, err.Error()))

		}
	})
	router.POST("/api/app/orc-server/build/file", func(context *gin.Context) {
		tenantId := context.GetHeader(consts.TENANT_ID)
		if file, err := context.FormFile("file"); err != nil {
			log.Error().Msgf("无法读取文件,key=file,%v", err)
			context.JSON(400, consts.NewBusinessException(1080500, err.Error()))
			return
		} else {
			savePath := path.Join(wd, "example", tenantId, file.Filename)

			if err = service.SaveUploadedFile(file, savePath); err != nil {
				log.Error().Msgf("无法保存文件,key=file,%v", err)
				context.JSON(400, consts.NewBusinessException(1080500, err.Error()))
				return
			}
			ch := make(chan error, 1)
			go func(channel chan error, filePath string) {
				//读取文件内容，并构建
				_, err1 := service.BuildJsonFile(savePath, tenantId)
				channel <- err1
			}(ch, savePath)
			select {
			case err = <-ch:
				if err != nil {
					context.JSON(400, consts.NewBusinessException(1080500, err.Error()))
					return
				}
			case <-time.After(5 * time.Minute):
				context.JSON(400, consts.NewBusinessException(1080500, "构建超时请检查"))
				return
			}
			context.JSON(200, consts.NewBusinessException(0, "构建成功"))
		}
	})
	router.DELETE("/api/app/orc-server/build/file", func(context *gin.Context) {
		tenantId := context.GetHeader(consts.TENANT_ID)
		if path, ok := context.GetQuery("api"); !ok {
			context.JSON(http.StatusBadRequest, consts.NewBusinessException(1080500, "缺少query必填参数api,示例：?api=/api/app/orc/bw/edit"))
			return
		} else if method, ok := context.GetQuery("method"); !ok {
			context.JSON(http.StatusBadRequest, consts.NewBusinessException(1080500, "缺少query必填参数api,示例：?method=post"))
			return
		} else {
			json, _ := context.GetQuery("json")
			context.JSON(service.Remove(tenantId, path, method, json))
			return
		}
	})
	router.POST("/api/app/orc-server/runner", func(context *gin.Context) {
		context.JSON(service.Runner(context))
		return
	})
	router.Any(*relativePath+"*action", func(context *gin.Context) {
		ch := make(chan error, 1)
		var result any
		go func(channel chan error, ctx *gin.Context) {
			//获取请求体
			//r, e := bweditpost.Executorbweditpost(ctx)
			r, e := service.Execute(ctx)
			channel <- e
			result = r
		}(ch, context)
		select {
		case err := <-ch:
			if err != nil {
				context.JSON(400, consts.NewBusinessException(1080500, err.Error()))
				return
			}
		case <-time.After(30 * time.Second):
			context.JSON(400, consts.NewBusinessException(1080500, "请求超时请检查"))
			return
		}
		context.JSON(http.StatusOK, consts.NewBusinessExceptionWithData(0, "请求成功", result))

	})

	port := strconv.Itoa(*serverPort)
	router.Run(":" + port)

}
