package js

import (
	"github.com/dop251/goja"
	"github.com/rs/zerolog/log"
	"strings"
	"testing"
)

func TestExecuteJavaScript(t *testing.T) {
	var script = "let records = $23f398a93f764236b0c2b7a05d6a1feb.$resp.data.data.records\nif(records.length === 0){\n    null\n}else{\n    let ok = {\n        \"blackGuid\":records[0].factoryRecordIdentify,\n        \"factoryIdentify\":records[0].factoryIdentify\n    }\n    ok\n}"

	split := strings.Split(script, "\n")

	for _, s := range split {
		println(s)
	}
}

type async struct {
	Vm   *goja.Runtime
	Func string
}

func (a *async) Suspended() (trackingObject interface{}) {
	log.Info().Msgf("调用了async function，我被调用了")
	if v, e := a.Vm.RunString(a.Func); e != nil {
		return e
	} else {
		return v.Export()
	}
}
func (a *async) Resumed(trackingObject interface{}) {
	log.Info().Msgf("我被调用了，但我不知道我是啥")
}

func TestAsyncJavaScript(t *testing.T) {
	script := `
let a = false
async function sleep(time) {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve()
        }, time)

    })
}

async function f() {
    //等待
    await sleep(2000);
    a = true
}
f()

`
	vm := goja.New()
	//	a := &async{vm, `function f() {
	// a = true
	//}`}
	//	vm.SetAsyncContextTracker(a)
	if v, e := vm.RunString(script); e != nil {
		t.Fatal("js执行异常", e)
	} else {
		promise := v.Export().(*goja.Promise)
		res := promise.Result()
		switch s := promise.State(); s {
		case goja.PromiseStateFulfilled:
			t.Logf("%+v", res)
		case goja.PromiseStateRejected:
			if resObj, ok := res.(*goja.Object); ok {
				if stack := resObj.Get("stack"); stack != nil {
					t.Error(stack.String())
				}
			}
		default:
			t.Fatalf("Unexpected promise state: %v", s)
		}
	}
}
