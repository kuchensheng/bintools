package http

import (
	"encoding/json"
	"fmt"
	"github.com/kuchensheng/bintools/logger"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"
)

type testError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e testError) Error() string {
	return e.Message
}

func TestEngine_Get(t *testing.T) {
	e := Default()
	e.Use(func(ctx *Context) {
		t.Logf("my = %s", "库陈胜")
		ctx.Set("my", "库陈胜")
		ctx.Next()
		name := ctx.GetString("name")
		logger := ctx.Logger()
		logger.Info("你好,%s", name)
		t.Logf("name = %s", name)
	})

	e.Get("/api/test", func(ctx *Context) {
		a, _ := ctx.GetQuery("a")
		t.Logf("a = %s", a)
		ctx.Set("name", "酷达舒")
		t.Logf("b= %s", ctx.GetString("b"))
		ctx.JSON(200, testError{
			Code:    0,
			Message: "成功",
			Data: struct {
				Name string `json:"name"`
				Ages []int  `json:"ages"`
			}{"嘿嘿", []int{1, 2, 3}},
		})
	})

	//从这以后的请求，将执行两个middleware
	e.Use(func(ctx *Context) {
		ctx.Set("b", "myb")
		//ctx.Next()
		name := ctx.GetString("name")
		t.Logf("bname = %s", name)
	})

	e.Get("/api/test/:id", func(ctx *Context) {
		id, _ := ctx.GetPath("id")
		t.Logf("aid = %s", id)
		ctx.Set("name", "酷达舒")
		t.Logf("b= %s", ctx.GetString("b"))
		ctx.JSON(200, testError{
			Code:    0,
			Message: "成功",
			Data: struct {
				Name string `json:"name"`
				Ages []int  `json:"ages"`
			}{"嘿嘿", []int{1, 2, 3, 4}},
		})
	})

	e.Get("/api/test/:id/:name", func(ctx *Context) {
		a, _ := ctx.GetPath("id")
		b, _ := ctx.GetPath("name")
		t.Logf("aid = %s", a)
		t.Logf("bName = %s", b)
		ctx.Set("name", "酷达舒")
		t.Logf("b= %s", ctx.GetString("b"))
		//go func(context *Context) {
		//	//todo 显示传递
		//
		//}(ctx)
		println("协程数量 = ", runtime.NumGoroutine())
		for i := 0; i < 100; i++ {
			go func(idx int) {
				time.Sleep(time.Second)
				println("id =", idx, "\t协程数量 = ", runtime.NumGoroutine())
				//隐式获取context

			}(i)
		}

		panic("我错误了")

		ctx.JSON(200, testError{
			Code:    0,
			Message: "成功",
			Data: struct {
				Name string `json:"name"`
				Ages []int  `json:"ages"`
			}{"嘿嘿", []int{1, 2, 3, 4, 5}},
		})
	})
	e.Any("/test", func(ctx *Context) {
		ctx.JSON(200, testError{
			Code:    400,
			Message: "我错了，打我呀",
		})
	})
	e.Any("/test/test1/*action", func(ctx *Context) {
		ctx.JSON(200, testError{
			Code:    400,
			Message: "我错了，打我呀111",
		})
	})

	e.RunWithPort(8080)
}

func TestEngine_PostForm(t *testing.T) {
	e := Default()
	e.Post("/api/test/post", func(ctx *Context) {
		form := ctx.PostForm("name")
		ctx.JSONoK(Result{0, "成功", form})
	})
	e.RunWithPort(8080)
}

func TestEngine_PostFormFile(t *testing.T) {
	e := Default()
	e.Post("/api/test/post", func(ctx *Context) {
		form, _ := ctx.FormFile("file")
		ctx.JSONoK(Result{0, "成功", form})
	})
	e.RunWithPort(8080)
}

func TestEngine_PostRawData(t *testing.T) {
	e := Default()
	e.Post("/api/test/post", func(ctx *Context) {
		data, _ := ctx.GetRawData()
		var a any
		json.Unmarshal(data, &a)
		ctx.JSONoK(Result{0, "成功", a})
	})
	e.RunWithPort(8080)
}

func TestEngine_GetWithUse(t *testing.T) {
	e := Default()
	//e.Pprof(true)
	e.Use(func(ctx *Context) {
		t.Logf("my = %s", "库陈胜")
		ctx.Set("my", "库陈胜")
		ctx.Next()
		name := ctx.GetString("name")
		t.Logf("name = %s", name)
	}, func(ctx *Context) {
		t.Logf("你好好:%s", "csd")
		ctx.Next()
		t.Log("我要稍后执行")
	})
	e.Get("/api/test", func(ctx *Context) {
		a, _ := ctx.GetQuery("a")
		t.Logf("a = %s", a)
		ctx.Set("name", "酷达舒")
		b, _ := ctx.GetQuery("b")
		t.Logf("b= %s", b)
		ctx.JSON(200, testError{
			Code:    0,
			Message: "成功",
			Data: struct {
				Name string `json:"name"`
				Ages []int  `json:"ages"`
			}{"嘿嘿", []int{1, 2, 3}},
		})
	})
	e.RunWithPort(8080)
}

func BenchmarkEngine_Get(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func(idx int) {
			if r, e := http.Get("http://localhost:8080/api/test?a=库陈胜ccc&b=帅不帅"); e != nil {
				b.Log(e)
			} else {
				data, _ := ioutil.ReadAll(r.Body)
				var result interface{}
				_ = json.Unmarshal(data, &result)
				b.Log(result)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestEngine_GetWithParam(t *testing.T) {
	e := Default()
	e.GetWithParam("/api/test/param", func(params ...HandlerParam) (any, error) {
		msg := fmt.Sprintf("param name = %s,value = %s", params[0].Name(), params[0].Value())
		logger.GlobalLogger.Info(msg)
		return msg, nil
	}, NewQuery("name", false), QueryParam{"age", 0, true}, BodyParam{struct {
		Class string
	}{""}, false})
	//e.AnyWithParam("/api/test/param", QueryParam{""}, QueryParam{0}, RequestBody{[]string{}})
	e.RunWithPort(8080)
}

type myBody struct {
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Datas []string `json:"datas"`
}

func TestEngine_PostWithParam(t *testing.T) {
	e := Default()
	e.PostWithParam("/api/test/param", func(params ...HandlerParam) (any, error) {
		return fmt.Sprintf("%+v", params), nil
	}, NewQuery("name", false), BodyParam{
		myBody{}, true,
	})
	e.RunWithPort(8080)
}

func TestEngine_Delete(t *testing.T) {
	e := Default()
	e.DeleteWithParam("/api/test/param", func(params ...HandlerParam) (any, error) {
		return fmt.Sprintf("%+v", params), nil
	}, NewQuery("name", false), BodyParam{
		myBody{}, true,
	})
	e.RunWithPort(8080)
}
