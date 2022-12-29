package log

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var channelMap = cache.New(time.Minute, 10*time.Second)
var END = "EOF"
var format = "2006-01-02 15:04:05.999"

type LogStruct struct {
	PK      string
	TraceId string
}

func (p LogStruct) Info(msg string, args ...any) {
	//日志格式:$DATE $TIME [funcName] [traceId] msg
	push(p.PK, "info", msg, p, buildMsg, args)
}

func (p LogStruct) Warn(msg string, args ...any) {
	//日志格式:$DATE $TIME [funcName] [traceId]
	push(p.PK, "warn", msg, p, buildMsg, args)
}

func (p LogStruct) Error(msg string, args ...any) {
	//日志格式:$DATE $TIME [funcName] [traceId] msg
	push(p.PK, "error", msg, p, buildMsg, args)
}

func buildMsg(level, msg string, p LogStruct, args ...any) string {
	_, f, line, _ := runtime.Caller(3)
	if !(len(args) == 1 && reflect.ValueOf(args[0].([]any)[0]).Len() == 0) {
		msg = fmt.Sprintf(msg, args[0].([]any)[0].([]any)...)
	}
	return fmt.Sprintf("%s [%s] %s:%d [%s] %s", now(), strings.ToUpper(level), f, line, p.TraceId, msg)
}

func now() string {
	return time.Now().Format(format)
}
func StartListener(pk string) {
	ch := make(chan string, 128)
	ch <- "连接成功"
	channelMap.SetDefault(pk, ch)
}
func StopListener(pk string) {
	channelMap.Delete(pk)
}
func push(pk string, level, msg string, p LogStruct, data func(level, msg string, p LogStruct, args ...any) string, args ...any) {
	if c, ok := channelMap.Get(pk); ok {
		c.(chan string) <- data(level, msg, p, args)
	} else {
		//未初始化，不执行push操作
	}
}

func Pull(pk string, conn *websocket.Conn) string {
	defer func() {
		conn.Close()
	}()
	if c, ok := channelMap.Get(pk); ok {
		for {
			select {
			case value := <-c.(chan string):
				conn.WriteMessage(websocket.TextMessage, []byte(value))
			case <-time.After(5 * time.Second):
				log.Warn().Msgf("5s内未取到值,结束监听")
				conn.WriteMessage(websocket.TextMessage, []byte(END))
				return END
			}
		}

	}
	return END
}
