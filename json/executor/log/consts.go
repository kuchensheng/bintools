package log

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"strings"
)

func GetPackage(ctx *gin.Context) string {
	method := strings.ToLower(ctx.Request.Method)
	version := ctx.GetHeader("version")
	return GetKey(ctx.Request.URL.Path, method, version)
}

func GetKey(uri, method, version string) string {
	key := strings.Join([]string{uri, method, version}, "")
	key = strings.ReplaceAll(key, consts.GlobalPrefix, "")
	key = strings.ReplaceAll(key, consts.GlobalTestPrefix, "")
	key = strings.ReplaceAll(key, "/", "")
	key = strings.ReplaceAll(key, "-", "")
	return key
}
