package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/configuration"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/lib"
	"github.com/kuchensheng/bintools/json/middleware"
	"github.com/kuchensheng/bintools/json/register"
	"github.com/kuchensheng/bintools/json/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func main() {
	//启动http服务，承接请求
	log.Logger = log.Logger.Level(zerolog.InfoLevel)
	wd, _ := os.Getwd()
	router := gin.Default()
	loginEnabled := false
	if v := configuration.GetConfig("login.enable"); v != nil {
		loginEnabled = v.(bool)
	}
	router.Use(middleware.TracerFilter())
	if loginEnabled {
		router.Use(middleware.LoginFilter())
	}

	router.POST("/api/app/orc-server/build", func(context *gin.Context) {
		data, _ := context.GetRawData()
		tenantId := context.GetHeader(consts.TENANT_ID)
		if _, err := service.BuildJson(data, tenantId); err != nil {
			context.JSON(400, consts.NewBusinessException(1080500, err.Error()))

		}
	})
	router.POST("/api/app/orc-server/build/file", func(context *gin.Context) {
		var tenantId string
		if t, ok := context.Get(consts.TENANT_ID); ok {
			tenantId = fmt.Sprintf("%v", t)
		}
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

	relativePath := consts.GlobalPrefix
	if v := configuration.GetConfig("server.context"); v != nil {
		relativePath = fmt.Sprintf("%v", v)
		consts.GlobalPrefix = relativePath
	}
	testRelativePath := consts.GlobalTestPrefix
	if v := configuration.GetConfig("server.context-test"); v != nil {
		testRelativePath = fmt.Sprintf("%v", v)
		consts.GlobalTestPrefix = testRelativePath
	}

	if !strings.HasSuffix(relativePath, "/") {
		relativePath += "/"
	}
	goPath := lib.GoPath
	if v := configuration.GetConfig("server.go.path"); v != nil {
		goPath = fmt.Sprintf("%v", v)
		lib.GoPath = goPath
	}
	routeHost := register.RouteHost
	if v := configuration.GetConfig("server.route.host"); v != nil {
		routeHost = fmt.Sprintf("%v", v)
		register.RouteHost = routeHost
	}

	go register.InitRoute()

	router.Any(relativePath+"*action", service.Execute)
	router.Any("/api/app/test/orc/*action", service.Execute)
	router.Any("/ws/app/orc/log", service.LogServer)

	//port := strconv.Itoa(*serverPort)
	port := "38240"
	if v := configuration.GetConfig("server.port"); v != nil {
		port = strconv.Itoa(v.(int))
	}
	router.Run(":" + port)

}
