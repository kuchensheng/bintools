package log

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetPackage(ctx *gin.Context) string {
	uri := strings.ReplaceAll(ctx.Request.URL.Path, "/", "_")
	method := strings.ToLower(ctx.Request.Method)
	version := ctx.GetHeader("version")
	return GetKey(uri, method, version)
}

func GetKey(uri, method, version string) string {
	key := strings.Join([]string{uri, method, version}, "_")
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "_api_app_orc_", "")
	key = strings.ReplaceAll(key, "_", "")
	return key
}
