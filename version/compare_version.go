package main

import (
	"fmt"
	"os"
	"strings"
)

//比较版本号大小
const (
	VersionBig   = 1
	VersionSmall = 2
	VersionEqual = 3
)

var splitFlag = "."

func main() {

	args := os.Args
	versionA := args[1]
	versionB := args[2]
	if len(args) > 2 {
		splitFlag = args[3]
	}

	fmt.Sprintln("compare version:", versionA, versionB, compareVersion(versionA, versionB))
}

func compareVersion(versionA, versionB string) int {
	fmt.Sprintf("compare version :%s,%s\n", versionA, versionB)
	//字符串分割
	versionAs := strings.Split(versionA, splitFlag)
	versionBs := strings.Split(versionB, splitFlag)

	lenVersionA := len(versionAs)
	lenVersionB := len(versionBs)

	if lenVersionA > lenVersionB {
		return VersionBig
	} else if lenVersionA < lenVersionB {
		return VersionSmall
	} else {
		return compareVersionItems(versionAs, versionBs)
	}
	return VersionEqual
}

//compareVersionItems 倒序，依次比较
func compareVersionItems(arrA, arrB []string) int {
	for i := 0; i < len(arrA); i++ {
		if littleResult := compareLittleVersion(arrA[i], arrB[i]); littleResult != VersionEqual {
			return littleResult
		}
	}
	return VersionEqual
}

func compareLittleVersion(itemA, itemB string) int {
	bytesA := []byte(itemA)
	bytesB := []byte(itemB)

	lenA := len(bytesA)
	lenB := len(bytesB)

	if lenA > lenB {
		return VersionSmall
	}
	//如果长途相等，则按byte位比较
	return func(verA, verB []byte) int {
		for index, _ := range verA {
			if verA[index] > verB[index] {
				return VersionBig
			}
			if verA[index] != verB[index] {
				return VersionSmall
			}
		}
		return VersionEqual
	}(bytesA, bytesB)
}
