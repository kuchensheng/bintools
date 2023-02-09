package logger

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

func (l Logger) enabled(lvl Level) bool {
	return lvl >= l.level
}

func (l *Logger) msg(lvl Level, msgFmt string, args ...any) string {
	//todo 字符串拼接而非直接创建字符串
	caller := l.formatter.caller(l.callerSkip)
	time := l.formatter.timeFmt(l.FormatTime)
	level := l.formatter.levelFmt(lvl.GetName())
	var items []string
	items = append(items, time, hostName, l.appName, l.traceId(), level, caller, l.dict(), fmt.Sprintf(msgFmt, args...), "\n")
	return strings.Join(items, " ")
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

func (l *Logger) Info(format string, args ...any) {
	if l.enabled(InfoLevel) {
		msg := l.msg(InfoLevel, format, args...)
		_ = l.WriteLevel(InfoLevel, msg)
	}
}

func (l Logger) Warn(format string, args ...any) {
	if l.enabled(WarnLevel) {
		format += "\n%-6s"
		args = append(args, debug.Stack())
		msg := l.msg(WarnLevel, format, args...)
		_ = l.WriteLevel(WarnLevel, msg)
	}
}

func (l Logger) Error(format string, args ...any) {
	if l.enabled(ErrorLevel) {
		format += "\n%-6s"
		args = append(args, debug.Stack())
		msg := l.msg(ErrorLevel, format, args...)
		_ = l.WriteLevel(ErrorLevel, msg)
	}
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

func (l Logger) FatalLevel(info any) {
	msg := ""
	if v, ok := info.(error); ok {
		msg = v.Error()
	} else {
		msg = fmt.Sprintf("%v", info)
	}
	_ = l.WriteLevel(FatalLevel, msg)
	os.Exit(0)
}
