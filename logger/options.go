package logger

import (
	"fmt"
	"os"
	"runtime/debug"
)

func (l Logger) enabled(lvl Level) bool {
	return lvl >= l.level
}

func (l Logger) Trace(format string, args ...any) {
	if l.enabled(TraceLevel) {
		_ = l.WriteLevel(TraceLevel, format, args...)
	}
}

func (l Logger) Debug(format string, args ...any) {
	if l.enabled(DebugLevel) {
		_ = l.WriteLevel(DebugLevel, format, args...)
	}
}

func (l Logger) Info(format string, args ...any) {
	if l.enabled(InfoLevel) {
		_ = l.WriteLevel(InfoLevel, format, args...)
	}
}

func (l Logger) Warn(format string, args ...any) {
	if l.enabled(WarnLevel) {
		format += "\n%-6s"
		args = append(args, debug.Stack())
		_ = l.WriteLevel(WarnLevel, format, args...)
	}
}

func (l Logger) Error(format string, args ...any) {
	if l.enabled(ErrorLevel) {
		format += "\n%-6s"
		args = append(args, debug.Stack())
		_ = l.WriteLevel(ErrorLevel, format, args...)
	}
}

func (l Logger) Panic(info any) {
	msg := ""
	if v, ok := info.(error); ok {
		msg = v.Error()
	} else {
		msg = fmt.Sprintf("%v", info)
	}
	_ = l.WriteLevel(PanicLevel, "%s\n", msg)
	panic(info)
}

func (l Logger) FatalLevel(info any) {
	msg := ""
	if v, ok := info.(error); ok {
		msg = v.Error()
	} else {
		msg = fmt.Sprintf("%v", info)
	}
	_ = l.WriteLevel(FatalLevel, "%s\n%s", msg, debug.Stack())
	os.Exit(0)
}
