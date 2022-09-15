package trace

import (
	"encoding/json"
	"errors"
	"github.com/kuchensheng/bintools/tracer/conf"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestNewServerTracerWithHttpServer(t *testing.T) {
	handlerFunc := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//业务逻辑
	})
	http.Handle("/", handlerFunc)
	svr := http.Server{
		Addr: ":8080",
		Handler: ServerTraceHandler(func(writer http.ResponseWriter, request *http.Request) {
			//业务逻辑
		}),
	}
	svr.ListenAndServe()
}

func TestNewServerTracer(t *testing.T) {
	req := &http.Request{
		Header: map[string][]string{"Content-Type": {"application/json"}, "token": {"i am authorization info"}},
		URL: func() *url.URL {
			url, _ := url.Parse("http://localhost:8080?id=23")
			return url
		}(),
		Method: http.MethodGet,
	}
	//create a tracer of server
	serverTracer := NewServerTracer(req)
	//todo do business
	t.Logf("业务处理中，请稍后……")
	time.Sleep(time.Second * 2)
	//end trace after done business.if it is OK,call serverTracer.EndTraceOK(),else call serverTracer.EndTracerError().or you can call serverTracer.EndTracer(OK,"this is message")
	serverTracer.EndTraceOk()
	//it is error
	//serverTracer.EndTraceError(errors.New("there is error message"))
	//it is other
	//serverTracer.EndTrace(WARNING,"this is waring message")
}

func TestNew(t *testing.T) {
	type args struct {
		req *http.Request
	}
	requestArgs := func() *http.Request {
		request := &http.Request{
			Header: map[string][]string{},
		}
		request.URL, _ = url.Parse("http://localhost:8080?id=23")
		request.PostForm = make(map[string][]string)
		request.PostForm.Set("isyscoreOS", "3.1.0")
		//request.MultipartForm.Value["content"] = []string{"6666"}
		//request.MultipartForm.File["file"] = []*multipart.FileHeader{
		//	{
		//		Filename: "haha.txt",
		//		Size:     222,
		//	},
		//}
		//request.Body.Read([]byte("kucs is a lucky boy"))
		return request
	}()
	tests := []struct {
		name string
		args args
		want *Tracer
	}{

		{
			name: "",
			args: struct{ req *http.Request }{
				req: requestArgs,
			},
			want: New(requestArgs),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				data, _ := json.Marshal(got)
				t.Logf("New() = %v, want %v", string(data), tt.want)
			}
		})
	}
}

func TestTracer_EndTraceOk(t *testing.T) {
	requestArgs := func() *http.Request {
		request := &http.Request{
			Header: func() map[string][]string {
				headers := make(map[string][]string)
				headers[T_HEADER_RPCID] = []string{"1.1"}
				return headers
			}(),
		}
		request.URL, _ = url.Parse("http://localhost:8080?id=23")
		request.Method = "GET"
		request.PostForm = make(map[string][]string)
		request.PostForm.Set("isyscoreOS", "3.1.0")
		return request
	}()

	tests := []struct {
		name string
		req  *http.Request
	}{
		{
			name: "测试01",
			req:  requestArgs,
		},
	}
	header := &http.Header{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf.Conf.Loki.Host = "http://10.30.30.78:3100"
			conf.Conf.Loki.MaxWaitTime = 1
			serverTracer := NewServerTracer(tt.req)
			//serverTracer1 := NewServerTracerWithoutReq()
			println("服务端其他业务请求")
			println("向客户端发起请求")
			for i := 0; i < 3; i++ {
				//clientTracer := serverTracer.NewClientTracer(tt.req)
				clientTracer := serverTracer.NewClientWithHeader(header)
				clientTracer.TraceName = "自定义traceName，默认:<Method>uri"
				clientTracer.AttrMap = []Parameter{}
				println("真正的请求，dorequest")
				//请求结束后，调用Endtrace
				clientTracer.EndTrace(OK, "i am danger")
			}
			//服务端请求结束后，调用EndTrace()
			//serverTracer.EndTrace(OK, "i am not in danger")
			err := errors.New("我打江南走过，大哥，我错了")
			serverTracer.EndTrace(ERROR, err.Error())
			time.Sleep(2 * time.Second)
		})
	}
}
