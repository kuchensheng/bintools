package util

import (
	"reflect"
)

func MakeStruct(names []string, vals []any) reflect.Value {
	if len(names) != len(vals) {
		panic("array out of index")
	}
	valMap := func(ns []string, vs []any) map[string]any {
		vm := make(map[string]any, len(ns))
		for i, v := range vs {
			vm[ns[i]] = v
		}
		return vm
	}(names, vals)
	return MakeStructByMap(valMap)
}

//MakeStructByMap 根据给定的字段，创建结构体
func MakeStructByMap(vals map[string]any) reflect.Value {
	var sfs []reflect.StructField

	for k, v := range vals {
		t := reflect.TypeOf(v)
		sf := reflect.StructField{
			Name: k,
			Type: t,
		}
		sfs = append(sfs, sf)
	}

	st := reflect.StructOf(sfs)
	so := reflect.New(st)

	return so
}
