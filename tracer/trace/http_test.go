package trace

import (
	"github.com/kuchensheng/bintools/tracer/conf"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestServerTracer_Delete(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "删除接口测试",
			fields: fields{
				&Tracer{
					TraceId:     LocalIdCreate.GenerateTraceId(),
					sampled:     true,
					ServiceName: conf.Conf.ServiceName,
					startTime:   time.Now().UnixMilli(),
					RpcId:       "0",
					TraceType:   HTTP,
					RemoteIp:    GetLocalIp(),
					TraceName:   "<default>_server",
				},
				"",
			},
			args: args{
				url:          "http://10.30.30.78:38080/api/apix/execute",
				header:       map[string][]string{"id": {"kucs"}},
				parameterMap: map[string]string{"name": "库陈胜"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.Delete(tt.args.url, tt.args.header, tt.args.parameterMap)
			t.Logf("clientId=%s", server.clientRpcId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_DeleteOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.DeleteOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_DeleteSimple(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.DeleteSimple(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_DeleteSimpleOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.DeleteSimpleOfStandard(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_Get(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.Get(tt.args.url, tt.args.header, tt.args.parameterMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_GetOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.GetOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_GetSimple(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.GetSimple(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_GetSimpleOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.GetSimpleOfStandard(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_Head(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			if err := server.Head(tt.args.url, tt.args.header, tt.args.parameterMap); (err != nil) != tt.wantErr {
				t.Errorf("Head() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerTracer_HeadSimple(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			if err := server.HeadSimple(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("HeadSimple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerTracer_Patch(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
		body         any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.Patch(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Patch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PatchOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
		body         any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.PatchOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PatchSimple(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url  string
		body any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.PatchSimple(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PatchSimpleOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url  string
		body any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.PatchSimpleOfStandard(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_Post(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
		body         any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.Post(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PostOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
		body         any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.PostOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PostSimple(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url  string
		body any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.PostSimple(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PostSimpleOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url  string
		body any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.PostSimpleOfStandard(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_Put(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
		body         any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.Put(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Put() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PutOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url          string
		header       http.Header
		parameterMap map[string]string
		body         any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.PutOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PutSimple(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url  string
		body any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.PutSimple(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PutSimpleOfStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		url  string
		body any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.PutSimpleOfStandard(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_call(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		httpRequest *http.Request
		url         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.call(tt.args.httpRequest, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("call() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_callIgnoreReturn(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		httpRequest *http.Request
		url         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			if err := server.callIgnoreReturn(tt.args.httpRequest, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("callIgnoreReturn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerTracer_callToStandard(t *testing.T) {
	type fields struct {
		Tracer      *Tracer
		clientRpcId string
	}
	type args struct {
		httpRequest *http.Request
		url         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &ServerTracer{
				Tracer:      tt.fields.Tracer,
				clientRpcId: tt.fields.clientRpcId,
			}
			got, err := server.callToStandard(tt.args.httpRequest, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("callToStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("callToStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetHttpClient(t *testing.T) {
	type args struct {
		httpClientOuter *http.Client
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetHttpClient(tt.args.httpClientOuter)
		})
	}
}

func Test_createHTTPClient(t *testing.T) {
	tests := []struct {
		name string
		want *http.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createHTTPClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseStandard(t *testing.T) {
	type args struct {
		responseResult []byte
		errs           error
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStandard(tt.args.responseResult, tt.args.errs)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_urlWithParameter(t *testing.T) {
	type args struct {
		url          string
		parameterMap map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := urlWithParameter(tt.args.url, tt.args.parameterMap); got != tt.want {
				t.Errorf("urlWithParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}
