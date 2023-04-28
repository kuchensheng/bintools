package util

import (
	"reflect"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	type args[T any] struct {
		array []T
		item  T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[string]{
		// TODO: Add test cases.
		{"addTest", args[string]{[]string{"1", "wo", "库"}, "陈"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.array, tt.args.item); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	type args[T any] struct {
		array []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[any]{
		// TODO: Add test cases.
		{"nilTest", args[any]{nil}, true},
		{"nilTest", args[any]{[]any{1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.array); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasNil(t *testing.T) {
	type args[T any] struct {
		array []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[any]{
		{"noNilTest", args[any]{[]any{1, "d", int32(1)}}, false},
		{"noNilTest", args[any]{[]any{1, "d", nil}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasNil(tt.args.array...); got != tt.want {
				t.Errorf("HasNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindFirst(t *testing.T) {
	type args[T any] struct {
		array     []T
		predicate func(item T) bool
	}
	type testCase[T any] struct {
		name  string
		args  args[T]
		want  T
		want1 bool
	}
	tests := []testCase[string]{
		{"aFirst", args[string]{array: []string{"你好", "导弹", "a", "cd"}, predicate: func(item string) bool {
			return item == "a"
		}}, "a", true},
		{"bFirst", args[string]{array: []string{"你好", "导弹", "a", "cd"}, predicate: func(item string) bool {
			return item == "b"
		}}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FindFirst(tt.args.array, tt.args.predicate)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindFirst() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindFirst() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAnyMatch(t *testing.T) {
	type args[T any] struct {
		array     []T
		predicate func(item T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[string]{
		// TODO: Add test cases.
		{"matchTest", args[string]{[]string{"1"}, func(item string) bool {
			return strings.Contains(item, "b")
		}}, false},
		{"matchTest", args[string]{[]string{"1", "bcdds", "cd"}, func(item string) bool {
			return strings.Contains(item, "b")
		}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyMatch(tt.args.array, tt.args.predicate); got != tt.want {
				t.Errorf("AnyMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Source struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Target struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Age   int    `json:"age"`
}

func TestMapper(t *testing.T) {
	type args[T any, R any] struct {
		array   []T
		convert func(item T) (R, error)
	}
	type testCase[T any, R any] struct {
		name    string
		args    args[T, R]
		want    []R
		wantErr bool
	}
	tests := []testCase[Source, Target]{
		// TODO: Add test cases.
		{"0Test", args[Source, Target]{array: []Source{{Name: "1", Age: 1}}, convert: func(item Source) (Target, error) {
			return Target{}, nil
		}}, []Target{{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mapper(tt.args.array, tt.args.convert)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mapper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mapper() got = %v, want %v", got, tt.want)
			}
		})
	}
}
