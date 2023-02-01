package http

import "testing"

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
		ctx.JSON(200, testError{
			Code:    0,
			Message: "成功",
			Data: struct {
				Name string `json:"name"`
				Ages []int  `json:"ages"`
			}{"嘿嘿", []int{1, 2, 3, 4, 5}},
		})
	})

	e.Run(8080)
}
