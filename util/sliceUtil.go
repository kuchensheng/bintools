package util

import (
	"reflect"
	"sort"
)

// IsEmpty 是否为空
func IsEmpty[T any](array []T) bool {
	return array == nil || len(array) == 0
}

// IsNotEmpty 是否非空
func IsNotEmpty[T any](array []T) bool {
	return !IsEmpty(array)
}

func DefaultIfEmpty[T any](array []T) []T {
	if IsEmpty(array) {
		array = make([]T, 0)
	}
	return array
}

// HasNil 是否包含nil元素
func HasNil[T any](array ...T) bool {
	if IsNotEmpty(array) {
		for _, item := range array {
			if reflect.DeepEqual(item, nil) {
				return true
			}
		}
	}
	return false
}

// FindFirst 返回第一个匹配到的值
func FindFirst[T any](array []T, predicate func(item T) bool) (T, bool) {
	for _, t := range array {
		if predicate(t) {
			return t, true
		}
	}
	return *new(T), false
}

// AnyMatch 判断是否存在匹配项
func AnyMatch[T any](array []T, predicate func(item T) bool) bool {
	source := array
	for len(source) > 0 {
		if len(source) == 1 {
			return predicate(source[0])
		}
		if predicate(source[0]) || predicate(source[len(source)-1]) {
			return true
		}

		source = source[1 : len(source)-1]
	}
	return false
}

func FirstMatchIndex[T any](array []T, predicate func(item T) bool) (int, bool) {
	for idx, item := range array {
		if predicate(item) {
			return idx, true
		}
	}
	return -1, false
}

// Mapper 将分片转成目标类型
func Mapper[T any, R any](array []T, convert func(item T) (R, error)) ([]R, error) {
	res := make([]R, len(array))
	for idx, t := range array {
		if r, err := convert(t); err != nil {
			return nil, err
		} else {
			res[idx] = r
		}
	}
	return res, nil
}

func Add[T any](array []T, item T) bool {
	length := len(array)
	array = append(array, item)
	return len(array) == length+1
}

func AddAll[T any](array []T, items []T) bool {
	length := len(array)
	array = append(array, items...)
	return len(array) == length+len(items)
}

func Insert[T any](array []T, item T, idx int) bool {
	length := len(array)
	if idx <= 0 {
		//插头
		array = append([]T{item}, array...)
	} else if idx >= length {
		//插尾
		array = append(array, item)
	} else {
		//插身
		left := append(array[:idx], item)
		right := array[idx:]
		array = append(left, right...)
	}
	return len(array) == length+1
}

func Set[T any](array []T, item T, idx int) {
	if idx >= 0 && idx < len(array) {
		array[idx] = item
	}
}

func Sort[T any](array []T, isAsc bool, compare func(o1, o2 T) bool) {
	sort.Slice(array, func(i, j int) bool {
		return compare(array[i], array[j]) && isAsc
	})
}

func Clear[T any](array []T) {
	array = nil
}
