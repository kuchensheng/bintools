package register

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

var RouteHost = "http://isc-route-service:31000"
var registerUri = "/api/route/refreshRoute"
var contentType = "application/json"

func InitRoute() {
	url := fmt.Sprintf("%s%s", RouteHost, registerUri)
	body := `{
	"path":"/api/app/orc-server/**;/api/app/orc/**",
	"serviceId":"isc-json-engine",
	"url":"http://isc-json-engine:38240"
}`
	if resp, err := http.Post(url, contentType, bytes.NewBufferString(body)); err != nil {
		log.Error().Msgf("路由规则注册失败,%v", err)
	} else {
		data := func() []byte {
			respBody := resp.Body
			defer respBody.Close()
			d, _ := ioutil.ReadAll(respBody)
			return d
		}()
		if resp.StatusCode != http.StatusOK {
			log.Warn().Msgf("路由规则注册失败,响应码返回:%d,响应内容:%s", resp.StatusCode, data)
		} else {
			log.Info().Msgf("路由规则注册成功")
		}
	}
}
