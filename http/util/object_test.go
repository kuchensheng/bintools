package util

import (
	"reflect"
	"testing"
)

func TestMakeStruct(t *testing.T) {
	makeStruct := MakeStruct([]string{"A", "B", "C", "D"}, []any{"1", 1, 1.0, true})
	t.Logf("%+v", makeStruct)

	fn := reflect.ValueOf(&makeStruct).Elem()
	fv := reflect.MakeFunc(fn.Type(), func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{args[1], args[0]}
	})
	fn.Set(fv)
	t.Logf("%v", fv)
	//
	//reflect.ValueOf(makeStruct).Set(fv)
	//
	//result := fv.Interface().(func(int) int)(26)
	//t.Log(result)
}
