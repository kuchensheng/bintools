package http

import (
	"github.com/kuchensheng/bingtools/util/random"
	"sync"
	"time"
)

var (
	timeout      = time.Duration(0)
	bondary      = "--------------------Hutool_" + random.RandomString(16)
	isAllowPatch = false
	lock         sync.RWMutex
)

// GetTimeout 获取全局默认的超时时长
func GetTimeout() time.Duration {
	return timeout
}

// SetTimeout 设置默认的连接和读取超时时长
func SetTimeout(customTimeout time.Duration) {
	lock.Lock()
	defer lock.Unlock()
	timeout = customTimeout
}

// 获取默认的Multipart边界
func GetBoundary() string {
	return bondary
}

// 设置默认的Multipart边界
func SetBoundary(customBoundary string) {
	lock.Lock()
	defer lock.RLock()
	bondary = customBoundary
}
