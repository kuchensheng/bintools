package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/configuration"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type UserStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    data   `json:"data"`
}

type data struct {
	Token      string   `json:"token"`
	UserId     string   `json:"userId"`
	LoginName  string   `json:"loginName"`
	NickName   string   `json:"nickName"`
	Role       []string `json:"role"`
	RoleId     []string `json:"roleId"`
	TenantId   string   `json:"tenantId"`
	SuperAdmin bool     `json:"superAdmin"`
}

var statusUri = "/api/permission/auth/status"
var defaultHost = "http://isc-permission-service:32100"

func getUserStatus(token string) UserStatus {
	status := UserStatus{
		Data: data{},
	}
	defer func() UserStatus {
		if x := recover(); x != nil {
			log.Error().Msgf("获取用户数据异常，%v", x)
		}
		return status
	}()

	header := http.Header{
		"token": []string{token},
	}
	if v := configuration.GetConfig("login.permissionUrl"); v != nil {
		defaultHost = fmt.Sprintf("%v", v)
	}
	url := func(h string) string {
		h = h + statusUri
		h = strings.ReplaceAll(h, "//api", "/api")
		return h
	}(defaultHost)
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = header
	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("不能正确地获取到用户状态,code = %d", resp.StatusCode)
		return status
	}
	if err != nil {
		log.Error().Msgf("获取用户状态时发生了异常,%v", err)
		return status
	}
	result := func() []byte {
		defer resp.Body.Close()
		ss, _ := ioutil.ReadAll(resp.Body)
		return ss
	}()
	if err = json.Unmarshal(result, &status); err != nil {
		log.Error().Msgf("用户状态信息解析异常,%v", err)
		return status
	}
	return status
}

var excludeUrl = func() []string {
	var res []string
	if v := configuration.GetConfig("login.exclude"); v != nil {
		for _, a := range v.([]any) {
			res = append(res, a.(string))
		}
	}
	return res
}()

func match(uri string) bool {
	for _, s := range excludeUrl {
		if uri == s {
			return true
		}
	}
	return false
}

func LoginFilter() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !match(context.Request.URL.Path) {
			token := context.GetHeader("token")
			if token == "" {
				context.JSON(http.StatusUnauthorized, consts.NewBusinessException(1080401, "登录校验失败,token为空"))
				context.Abort()
				return
			}
			status := getUserStatus(token)
			if status.Data.TenantId == "" {
				context.JSON(http.StatusUnauthorized, consts.NewBusinessException(1080401, "租户信息校验失败"))
				context.Abort()
				return
			}
			context.Set(consts.TENANT_ID, status.Data.TenantId)
		}
		context.Next()
	}
}
