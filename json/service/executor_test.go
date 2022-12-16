package service

import "testing"

func TestCompile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "测试文件",
			args:    args{"D:\\worksapace\\go\\bintools\\json\\example\\hahaha\\bweditpost.go_"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Compile(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
