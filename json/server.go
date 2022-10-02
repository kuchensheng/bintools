package main

import (
	"flag"
	"github.com/gin-gonic/gin"
)

func main() {
	//启动http服务，承接请求
	relativePath := flag.String("context_path", "/api/app/orc", "请求路径前缀")

	router := gin.Default()
	router.Any(*relativePath, func(context *gin.Context) {
		//todo 执行插件
	})

}
