package js

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/executor/util"
	"github.com/kuchensheng/bintools/tracer/trace"
	"github.com/rs/zerolog/log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var scriptEnginFunc = func(context *gin.Context) *goja.Runtime {
	scriptEngine := goja.New()
	scriptEngine.Set("ctx", context)
	scriptEngine.Set("getValueByKey", func(ctx *gin.Context, key string) any {
		value := util.GetContextValue(ctx, key)
		log.Info().Msgf("获取键=%s的值：%v", key, value)
		if value != nil {
			return value
		}
		return ""
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
			fmt.Printf("%s\n", debug.Stack())
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
		if v != nil || v.ExportType() != nil {
			return v.Export(), nil
		}
		return v.String(), nil
	}
}

func replaceScript(script string) string {
	log.Info().Msgf("替换前的脚本内容:%s", script)
	split := strings.Split(script, "\n")
	var placeholder []consts.Pair[string, string]
	var noSpaceLines []string
	for _, s := range split {
		if s != "" && strings.TrimSpace(s) != "" {
			noSpaceLines = append(noSpaceLines, s)
		}
	}
	for i, s := range noSpaceLines {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "return") {
			sb := strings.Builder{}
			for _, c := range placeholder {
				sb.WriteString("\n")
				sb.WriteString(fmt.Sprintf(`setValueByKey(ctx,"%s",%v)`, strings.TrimSpace(c.Second), c.First))
			}
			split[i] = fmt.Sprintf("%s\n%s\n", s, sb.String())
			placeholder = nil
		}

		if validToken(s) {
			noSpaceLines[i], placeholder = replaceGetOrSetValue(s, placeholder)
		}
	}

	script = strings.Join(noSpaceLines, "\n")
	for _, c := range placeholder {
		script = strings.ReplaceAll(script, c.Second, c.First)
	}
	sb := strings.Builder{}
	sb.Write([]byte(script))
	for _, c := range placeholder {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`setValueByKey(ctx,"%s",%v)`, strings.TrimSpace(c.Second), c.First))
	}
	script = sb.String()

	log.Info().Msgf("替换后的脚本内容:%s", script)
	return script
}

func replaceGetOrSetValue(s string, placeholder []consts.Pair[string, string]) (string, []consts.Pair[string, string]) {
	if strings.Contains(s, "=") {
		keys := strings.Split(s, "=")
		first := strings.TrimSpace(keys[0])
		first, placeholder = replaceGetOrSetValue(first, placeholder)
		second := strings.TrimSpace(keys[1])
		second, placeholder = replaceGetOrSetValue(second, placeholder)
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
		return strings.Join(keys, "="), placeholder
	} else if strings.Contains(s, ":") {
		keys := strings.Split(s, ":")
		//first := strings.TrimSpace(keys[0])
		second := strings.TrimSpace(keys[1])
		containsComman := strings.Contains(second, ",")
		if containsComman {
			second = strings.ReplaceAll(second, ",", "")
		}

		//获取值
		if validToken(second) {
			keys[1] = fmt.Sprintf(`getValueByKey(ctx,"%s")`, second)
		}
		if containsComman {
			keys[1] = keys[1] + ","
		}
		return strings.Join(keys, ":"), placeholder
	}
	return s, placeholder
}

func validToken(content string) bool {
	return strings.Contains(content, consts.KEY_TOKEN) && strings.Contains(content, consts.KEY_REQ_CONNECTOR)
}
