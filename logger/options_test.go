package logger

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

var logger = GlobalLogger

func TestTrace(t *testing.T) {
	logger.Trace("%v", "我是trace")
}

func TestNewWithCtx(t *testing.T) {
	ctx := context.WithValue(context.TODO(), TRACEID, "我是traceId")
	l := logger.WithContext(ctx)
	l.Info("侬好")
}

func TestLogger_WithContext(t *testing.T) {
	logger = logger.WithContext(context.WithValue(context.TODO(), "my", "1-2"))
}

func TestLogger(t *testing.T) {
	now := time.Now().UnixMilli()
	counter := 3000
	var sw sync.WaitGroup
	sw.Add(counter)
	l := logger
	l.CallerSkip(3)
	for i := 0; i < counter; i++ {
		go func(idx int) {
			logger.Info("%s%s,%6d", "库陈胜", "帅", i)
			sw.Done()
		}(i)
	}
	sw.Wait()
	now1 := time.Now().UnixMilli()
	logger.Info("执行完毕,耗时：%d ms", now1-now)
}

func BenchmarkLogger_Info(b *testing.B) {
	//w, _ := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 755)
	//logger = logger.MultiWriter(w, os.Stdout)
	//logger.Info("%s%s", "库陈胜", "帅吗？")
	now := time.Now().UnixMilli()
	counter := b.N
	var sw sync.WaitGroup
	sw.Add(counter)
	l := logger
	l.CallerSkip(3)
	for i := 0; i < counter; i++ {
		go func(idx int) {
			l.Info("%s%s,%s", "库陈胜", "帅", time.Now())
			sw.Done()
		}(i)
	}
	sw.Wait()
	now1 := time.Now().UnixMilli()
	logger.Info("执行完毕,耗时：%d ms", now1-now)
}

func TestLogger_Info(t *testing.T) {
	logger.Info("%s", "我是Info")
}
func TestWarn(t *testing.T) {
	logger.Warn("%v", errors.New("你好"))
}

func TestWarn2(t *testing.T) {
	ctx := context.WithValue(context.TODO(), TRACEID, "我是traceId")
	l := logger.WithContext(ctx)
	l.Warn("%v", errors.New("哈函数"))
}

func TestLogger_Debug(t *testing.T) {
	logger.Debug("%v", "我是debug")
}

func TestLogger_Error(t *testing.T) {
	logger.Error("%v", errors.New("我错误了"))
}

func TestLogger_Panic(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Errorf("%v", x)
		}
	}()
	logger.Panic("错误信息")
}

func TestLogger_FatalLevel(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Errorf("%v", x)
		}
	}()
	logger.FatalLevel("啊哈")
	t.Logf("你好：%s", "我执行了吗？")
}

func TestLogger_Dict(t *testing.T) {
	logger.Dict("name", "kucs")
	logger.Info("%s", "你好")
}
