package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
)

func (l Logger) enabled(lvl Level) bool {
	return lvl >= l.level
}

const space = ' '

var bytesPool = sync.Pool{
	New: func() any {
		return bytes.Buffer{}
	},
}

func (l *Logger) msg(lvl Level, msgFmt string, args ...any) string {
	//todo 字符串拼接而非直接创建字符串
	bs := bytesPool.Get().(bytes.Buffer)
	defer func(buffer bytes.Buffer) {
		buffer.Reset()
		defer bytesPool.Put(buffer)
	}(bs)
	bs.Reset()
	bs.WriteString(l.formatter.timeFmt(l.FormatTime))
	bs.WriteRune(space)
	bs.WriteString(hostName)
	bs.WriteRune(space)
	bs.WriteString(l.appName)
	bs.WriteRune(space)
	bs.WriteString(l.traceId())
	bs.WriteString(l.formatter.levelFmt(lvl.GetName()))
	bs.WriteRune(space)
	bs.WriteString(l.formatter.caller(l.callerSkip))
	bs.WriteRune(space)
	bs.WriteString(l.dict())
	bs.WriteRune(space)
	bs.WriteString(fmt.Sprintf(msgFmt, args...))
	bs.WriteRune('\n')
	return bs.String()
}

func (l Logger) Trace(format string, args ...any) {
	if l.enabled(TraceLevel) {
		msg := l.msg(TraceLevel, format, args...)
		_ = l.WriteLevel(TraceLevel, msg)
	}
}

func (l Logger) Debug(format string, args ...any) {
	if l.enabled(DebugLevel) {
		msg := l.msg(DebugLevel, format, args...)
		_ = l.WriteLevel(DebugLevel, msg)
	}
}

func (l Logger) Info(format string, args ...any) {
	if l.enabled(InfoLevel) {
		msg := l.msg(InfoLevel, format, args...)
		_ = l.WriteLevel(InfoLevel, msg)
	}
}

func (l Logger) Warn(format string, args ...any) {
	if l.enabled(WarnLevel) {
		if l.stack {
			format += "\n%-6s"
			args = append(args, debug.Stack())
		}
		msg := l.msg(WarnLevel, format, args...)
		_ = l.WriteLevel(WarnLevel, msg)
	}
}

func (l Logger) Stack() Logger {
	l.stack = true
	return l
}

func (l Logger) Error(format string, args ...any) {
	if l.enabled(ErrorLevel) {
		if l.stack {
			format += "\n%-6s"
			args = append(args, debug.Stack())
		}
		msg := l.msg(ErrorLevel, format, args...)
		_ = l.WriteLevel(ErrorLevel, msg)
	}
}

func (l Logger) Panicf(format string, args ...any) {
	l.Panic(fmt.Sprintf(format, args...))
}

func (l Logger) Panic(info any) {
	msg := ""
	if v, ok := info.(error); ok {
		msg = v.Error()
	} else {
		msg = fmt.Sprintf("%v", info)
	}
	_ = l.WriteLevel(PanicLevel, msg)
	panic(info)
}

func (l Logger) Fatalf(formt string, args ...any) {
	l.Fatal(fmt.Sprintf(formt, args...))
}

func (l Logger) Fatal(info any) {
	msg := ""
	if v, ok := info.(error); ok {
		msg = v.Error()
	} else {
		msg = fmt.Sprintf("%v", info)
	}
	_ = l.WriteLevel(FatalLevel, msg)
	os.Exit(0)
}
