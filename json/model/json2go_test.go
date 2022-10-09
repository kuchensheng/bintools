package model

import (
	"testing"
)

func TestApixData_GenerateGo(t *testing.T) {
	path, err := GenerateFile2Go("C:\\Users\\admin\\Desktop\\dsl\\first.dsl")
	if err != nil {
		t.Logf("通过模板构建文件内容错误,%v", err)
		t.Fatalf("通过模板构建文件内容错误,%v", err)
	}
	t.Logf("插件文件路径:%s", path)
}

func TestBuild(t *testing.T) {
	path, err := Build("C:\\Users\\admin\\Desktop\\dsl\\first.dsl")
	if err != nil {
		t.Logf("dsl编译失败,%v", err)
		t.Fatalf("dsl编译失败,%v", err)
	}
	for s, info := range GetAllPlugins() {
		t.Logf("pluginName=[%s],info=%v", s, info)
	}
	t.Logf("go源码文件路径:%s", path)
}
