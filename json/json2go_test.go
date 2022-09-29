package json

import (
	"net/http"
	"testing"
)

func TestApixData_GenerateGo(t *testing.T) {
	data := &ApixData{}

	rule := data.Rule

	api := rule.Api
	key := api.Path + api.Method
	println("插件的key:", key)

	println(rule)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//todo 执行插件
	})
	http.ListenAndServe(":8080", nil)
}
