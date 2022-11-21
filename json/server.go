package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	//"github.com/kuchensheng/bintools/json/example"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const defaultMultipartMemory = 2 << 20 // 32 MB

func main() {
	//启动http服务，承接请求
	relativePath := flag.String("context_path", "/api/app/orc/*action", "请求路径前缀")
	serverPort := flag.Int("port", 38240, "服务器端口，默认:38240")
	flag.Parse()
	log.Logger = log.Logger.Level(zerolog.InfoLevel)
	wd, _ := os.Getwd()
	router := gin.Default()
	router.POST("/api/app/orc-server/build", func(context *gin.Context) {
		data, _ := context.GetRawData()
		if _, err := model.BuildJson(data); err != nil {
			context.JSON(400, model.NewBusinessException(1080500, err.Error()))
		}
	})
	router.POST("/api/app/orc-server/build/file", func(context *gin.Context) {
		if file, err := context.FormFile("file"); err != nil {
			log.Error().Msgf("无法读取文件,key=file,%v", err)
			context.JSON(400, model.NewBusinessException(1080500, err.Error()))
			return
		} else {
			savePath := path.Join(wd, "example", file.Filename)

			if err = model.SaveUploadedFile(file, savePath); err != nil {
				log.Error().Msgf("无法保存文件,key=file,%v", err)
				context.JSON(400, model.NewBusinessException(1080500, err.Error()))
				return
			}
			ch := make(chan error, 1)
			go func(channel chan error, filePath string) {
				//读取文件内容，并构建
				_, err1 := model.BuildJsonFile(savePath)
				channel <- err1
			}(ch, savePath)
			select {
			case err = <-ch:
				if err != nil {
					context.JSON(400, model.NewBusinessException(1080500, err.Error()))
					return
				}
			case <-time.After(30 * time.Second):
				context.JSON(400, model.NewBusinessException(1080500, "构建超时请检查"))
				return
			}
			context.JSON(200, model.NewBusinessException(0, "构建成功"))
		}
	})
	router.Any(*relativePath, func(context *gin.Context) {
		//example.Executor(context.Request, context.Writer)
		model.ExecutePlugin(getPluginKey(context.Request), context)
	})

	port := strconv.Itoa(*serverPort)
	router.Run(":" + port)

}

func getPluginKey(request *http.Request) string {
	path := request.URL.Path
	method := request.Method
	method = strings.ToLower(method)

	version := request.URL.Query().Get("version")
	code := request.URL.Query().Get("code")
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
