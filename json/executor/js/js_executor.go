package js

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/executor/util"
	"github.com/kuchensheng/bintools/tracer/trace"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

var scriptEnginFunc = func(context *gin.Context) *goja.Runtime {
	scriptEngine := goja.New()
	scriptEngine.Set("ctx", context)
	scriptEngine.Set("getValueByKey", func(ctx *gin.Context, key string) string {
		return util.GetContextValue(ctx, key).(string)
	})
	scriptEngine.Set("setValueByKey", func(ctx *gin.Context, key string, value any) {
		util.SetResultValue(ctx, key, value)
	})
	//设定最长执行时间：1分钟
	time.AfterFunc(time.Minute, func() {
		scriptEngine.Interrupt("timeout")
	})
	return scriptEngine
}

//ExecuteJavaScript 执行JS脚本,返回执行结果或者错误信息
func ExecuteJavaScript(ctx *gin.Context, script, name string) (any, error) {
	tracer, _ := ctx.Get(consts.TRACER)
	clientTracer := tracer.(*trace.ServerTracer).NewClientWithHeader(&ctx.Request.Header)
	clientTracer.TraceName = "执行脚本节点:" + name
	defer func() {
		if x := recover(); x != nil {
			log.Error().Msgf("JS脚本执行异常，panic is :%v", x)
			clientTracer.EndTraceError(x.(error))
		}
	}()
	//初始化JS引擎
	scriptEngine := scriptEnginFunc(ctx)
	script = replaceScript(script)
	if v, err := scriptEngine.RunString(script); err != nil {
		clientTracer.EndTraceError(err)
		return nil, err
	} else {
		clientTracer.EndTraceOk()
		return v.Export(), nil
	}
}

func replaceScript(script string) string {
	log.Info().Msgf("替换前的脚本内容:%s", script)
	split := strings.Split(script, "\n")
	var placeholder []consts.Pair[string, string]
	for i, s := range split {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "return") {
			sb := strings.Builder{}
			for _, c := range placeholder {
				sb.WriteString("\n")
				sb.WriteString(fmt.Sprintf(`setValueByKey(ctx,"%s",%v)`, c.Second, c.First))
			}
			split[i] = fmt.Sprintf("%s\n%s\n", s, sb.String())
			placeholder = nil
		}

		if validToken(s) {
			keys := strings.Split(s, "=")
			first := strings.TrimSpace(keys[0])

			second := strings.TrimSpace(keys[1])
			//获取值
			if validToken(second) {
				keys[1] = fmt.Sprintf(`getValueByKey(ctx,"%s")`, second)
			}
			//赋值动作
			if validToken(first) {
				random := "a" + strconv.FormatInt(time.Now().UnixMilli(), 10)
				placeholder = append(placeholder, consts.Pair[string, string]{random, keys[0]})
				keys[0] = random
				if !strings.HasPrefix(keys[0], "let") {
					keys[0] = "let " + keys[0]
				}
			}
			split[i] = strings.Join(keys, "=")

		}
	}

	script = strings.Join(split, "\n")
	for _, c := range placeholder {
		script = strings.ReplaceAll(script, c.Second, c.First)
	}
	sb := strings.Builder{}
	sb.Write([]byte(script))
	for _, c := range placeholder {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`setValueByKey(ctx,"%s",%v)`, c.Second, c.First))
	}
	script = sb.String()

	log.Info().Msgf("替换后的脚本内容:%s", script)
	return script
}

func validToken(content string) bool {
	return strings.Contains(content, consts.KEY_TOKEN) && strings.Contains(content, consts.KEY_REQ_CONNECTOR)
}
