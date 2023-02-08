package logger

import (
	"fmt"
	"github.com/kuchensheng/bintools/logger/rotate"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	timeLayout = "2006-01-02"
	suffix     = ".log"
	checkSpec  = "*/1 * * * * ?"
)

func init() {
	var logger = GlobalLogger
	loggerTask := rotate.New()

	changeWriter := func(old *FileLevelWriter) {
		//创建文件
		newFile := newLogFile(old.File)
		//指针切换
		_ = os.Symlink(newFile.Name(), old.link)
	}
	fileHandler := func() {
		//修改源
		for _, writer := range logger.writer.writers {
			if w, ok := writer.(*FileLevelWriter); ok {
				changeWriter(w)
			}
		}
	}
	GlobalLogger = logger
	loggerTask.AddFunc(logger.spec, fileHandler)
	loggerTask.AddFunc(checkSpec, func() {
		logger.Info("检查文件大小")
		for _, writer := range GlobalLogger.writer.writers {
			if fw, ok := writer.(*FileLevelWriter); ok {
				fi, _ := fw.Stat()
				if fi.Size() >= logger.splitSize {
					changeWriter(fw)
				}
			}
		}
	})

}

func newLogFile(old *os.File) *os.File {
	dst := time.Now().Format(timeLayout)
	dstLog := dst + suffix
	dir := filepath.Dir(old.Name())
	if files, err := ioutil.ReadDir(dir); err != nil {
		GlobalLogger.Error("无法打开根目录[%s],%v", dir, err)
		return nil
	} else {
		one := lastOneByTime(files, func(item fs.FileInfo) bool {
			return strings.HasPrefix(item.Name(), dst)
		})
		fn := func(dst string) *os.File {
			f, e := os.Create(dst)
			if e != nil {
				GlobalLogger.Error("无法创建日志文件,%s,%v", dst, e)
				return nil
			} else {
				return f
			}
		}
		if one == nil {
			//创建文件
			return fn(dstLog)
		} else {
			//文件 + 标号
			name := one.Name()
			split := strings.Split(name, ".")
			length := len(split)
			last := split[length-1]
			if last == suffix {
				name += ".1"
			} else if idx, convErr := strconv.Atoi(last); convErr != nil {
				name += ".1"
			} else {
				name += fmt.Sprintf(".%d", idx+1)
			}
			return fn(name)
		}
	}
}

func lastOneByTime(list []fs.FileInfo, predicate func(item fs.FileInfo) bool) fs.FileInfo {
	list = filter(list, predicate)
	if list == nil {
		return nil
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ModTime().After(list[j].ModTime())
	})
	length := len(list)
	return list[length-1]
}

func filter(list []fs.FileInfo, filter func(item fs.FileInfo) bool) []fs.FileInfo {
	var result []fs.FileInfo
	for _, info := range list {
		if filter(info) {
			result = append(result, info)
		}
	}
	return result
}
