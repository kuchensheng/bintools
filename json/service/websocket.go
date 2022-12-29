package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kuchensheng/bintools/json/consts"
	log2 "github.com/kuchensheng/bintools/json/executor/log"
	"github.com/rs/zerolog/log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func LogServer(context *gin.Context) {
	if upgrade, err := upgrader.Upgrade(context.Writer, context.Request, nil); err != nil {
		log.Warn().Msgf("无法建立websocket连接:%v", err)
		context.JSON(http.StatusBadRequest, consts.NewBusinessException(1080400, "无法建立websocket连接:"+err.Error()))
		return
	} else {
		defer upgrade.Close()
		if api, ok := context.GetQuery("api"); !ok {
			upgrade.WriteMessage(websocket.TextMessage, []byte("api参数不能为空"))
			return
		} else if method, ok := context.GetQuery("method"); !ok {
			upgrade.WriteMessage(websocket.TextMessage, []byte("method参数不能为空"))
			return
		} else {
			version, _ := context.GetQuery("version")
			pk := log2.GetKey(api, method, version)
			log2.StartListener(pk)
			//持续监听
			_ = log2.Pull(pk, upgrade)
			//移除缓存
			log2.StopListener(pk)
		}
	}
}
