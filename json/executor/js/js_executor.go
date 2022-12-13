package js

import (
	"github.com/dop251/goja"
	"github.com/kuchensheng/bintools/json/consts"
	"os"
	"path/filepath"
	"time"
)

var scriptEngine *goja.Runtime

//ExecuteJavaScript 执行JS脚本,返回执行结果或者错误信息
func ExecuteJavaScript(script, name string) (any, error) {
	getwd, _ := os.Getwd()
	fp := filepath.Join(getwd, "scripts", name+".js")
	//首先判断文件是否存在
	if e := write2Script(fp, []byte(script)); e != nil {
		return nil, e
	}
	//初始化JS引擎
	if scriptEngine == nil {
		scriptEngine = goja.New()
	}

	//设定最长执行时间：1分钟
	time.AfterFunc(time.Minute, func() {
		scriptEngine.Interrupt("timeout")
	})
	var program *goja.Program
	if p, ok := consts.Cache.Get(name); !ok {
		if p, e := goja.Compile(name, fp, false); e != nil {
			return nil, e
		} else {
			program = p
			consts.Cache.SetDefault(name, program)
		}
	} else {
		program = p.(*goja.Program)
	}

	if v, err := scriptEngine.RunProgram(program); err != nil {
		return nil, err
	} else {
		return []byte(v.String()), nil
	}
}

func write2Script(path string, content []byte) error {
	if _, e := os.Stat(path); e != nil {
		if os.IsNotExist(e) {
			return os.WriteFile(path, content, 0666)
		}
		return e
	}
	return nil
}
