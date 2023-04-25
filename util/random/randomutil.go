package random

import (
	"math/big"
	"math/rand"
	"reflect"
)

var (
	baseNumber     = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	baseChar       = []int32{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	baseByte       = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	baseCharNumber = append(baseChar, []int32{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}...)
)

const (
	EMPTY_STR = ""
)

// RandomInt 指定范围的随机数[0,1),如果limit<=0，则认为没有范围限定
func RandomInt(limit int) int {
	if limit <= 0 {
		return rand.Int()
	}
	return rand.Intn(limit)
}

// RandomInts 生成指定长度的int数组/切片
func RandomInts(count int) []int {
	var res []int
	for i := 0; i < count; i++ {
		res[i] = RandomInt(0)
	}
	return res
}

// RandomInt32 指定范围的随机数[0,1),如果limit<=0，则认为没有范围限定
func RandomInt32(limit int32) int32 {
	if limit <= 0 {
		return rand.Int31()
	}
	return rand.Int31n(limit)
}

// RandomInt32s 生成指定长度的int32数组/切片
func RandomInt32s(count int) []int32 {
	var res []int32
	for i := 0; i < count; i++ {
		res[i] = rand.Int31()
	}
	return res
}

// RandomInt64 指定范围的随机数[0,1),如果limit<=0，则认为没有范围限定
func RandomInt64(limit int64) int64 {
	if limit <= 0 {
		return rand.Int63()
	}
	return rand.Int63n(limit)
}

// RandomInt64s 生成指定长度的int64数组/切片
func RandomInt64s(count int) []int64 {
	var res []int64
	for i := 0; i < count; i++ {
		res[i] = rand.Int63()
	}
	return res
}

// RandomFloat32 生成float32随机数
func RandomFloat32() float32 {
	return rand.Float32()
}

// RandomFloat32s 生成float32随机数切片
func RandomFloat32s(count int) []float32 {
	var res []float32
	for i := 0; i < count; i++ {
		res[i] = RandomFloat32()
	}
	return res
}

// RandomFloat64 生成float64随机数
func RandomFloat64() float64 {
	return rand.Float64()
}

// RandomFloat64s 生成float64随机数切片
func RandomFloat64s(count int) []float64 {
	var res []float64
	for i := 0; i < count; i++ {
		res[i] = RandomFloat64()
	}
	return res
}

func RandomBigInt() *big.Int {
	return big.NewInt(RandomInt64(0))
}

func RandomBigFloat() *big.Float {
	return big.NewFloat(RandomFloat64())
}

func RandomBigRat() *big.Rat {
	return big.NewRat(RandomInt64(0), RandomInt64(0))
}

// RandomBytes 组成给定长的byte数组
func RandomBytes(count int) []byte {
	var res []byte
	for i := 0; i < count; i++ {
		res[i] = baseByte[RandomInt(26)]
	}
	return res
}

// RandomEle 从list中随机获取一个元素
func RandomEle[T any](list []T) T {
	return list[RandomInt(len(list))]
}

// RandomEleList 从list中随机提取元素并组成新的Slice，可能存在重复
func RandomEleList[T any](list []T, count int) []T {
	var res []T
	for i := 0; i < count; i++ {
		res[i] = RandomEle(list)
	}
	return res
}

// RandomEleSet 从list中随机提取元素并组成新Set,不重复
func RandomEleSet[T any](list []T, count int) []T {
	//分片去重
	list = Distinct(list)
	var res []T
	keySet := make(map[int]int, count)
	for len(keySet) < count {
		randomIdx := RandomInt(len(list))
		if _, ok := keySet[randomIdx]; !ok {
			res = append(res, list[randomIdx])
			keySet[randomIdx] = randomIdx
		}
	}
	return res
}

// Distinct 切片去重，并返回新的切片
func Distinct[T any](list []T) []T {
	item := list[0]
	res := []T{item}
	//去重
	for _, t := range list {
		if !reflect.DeepEqual(t, item) {
			res = append(res, t)
			item = t
		}
	}
	return res
}

// RandomString 获取一个随机字符串(只包含数字和字符)
func RandomString(length int) string {
	return RandomStringWithBase(baseChar, length)
}

func RandomStringWithBase(baseChar []rune, length int) string {
	if baseChar == nil || len(baseChar) == 0 {
		return EMPTY_STR
	}
	res := make([]rune, length)
	for i := 0; i < length; i++ {
		res[i] = RandomChar()
	}
	return string(res)
}

func RandomChar() rune {
	return RandomCharWithBase(baseChar)
}

func RandomCharWithBase(baseChar []rune) rune {
	return baseChar[RandomInt(len(baseChar))]
}
