package yaml

import (
	"reflect"
	"testing"
)

func Test_readYaml(t *testing.T) {
	type args struct {
		yamlPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{"readYaml", args{"D:\\worksapace\\go\\bintools\\http\\config\\resources\\application.yml"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readYaml(tt.args.yamlPath)
		})
	}
}

func TestYamlConfig_GetAttr(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{"test", args{"server.port"}, 8080},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &YamlConfig{}
			if got := c.GetAttr(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Server struct {
	Port int
}
