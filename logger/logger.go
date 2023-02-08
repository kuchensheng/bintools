package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Level int8

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
	NoLevel
	Disabled
)

func (lvl Level) GetName() string {
	switch lvl {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		return ""
	}
}

var (
	hostName = func() string {
		name, _ := os.Hostname()
		return name
	}()
	pwd = func() string {
		path, _ := os.Getwd()
		return path
	}
	FmtYmdHmsSSS = "2006-01-02 15:04:05,000"
	TRACEID      = "t-head-traceId"
)

type CallerFun func(skip int) string
type FormatLevelFun func(i any) string
type LogFormatter struct {
	caller   CallerFun
	timeFmt  func(layout string) string
	levelFmt FormatLevelFun
}

type Logger struct {
	writer      multiWriter
	ctx         context.Context
	level       Level
	formatter   LogFormatter
	callerSkip  int
	appName     string
	FormatTime  string
	FormatLevel FormatLevelFun
	buffer      sync.Pool
	Keys        map[string]any
}

func New() Logger {
	l := Logger{
		level:      InfoLevel,
		callerSkip: 4,
		formatter: LogFormatter{
			caller: func(skip int) string {
				_, file, line, _ := runtime.Caller(skip)
				//当前路径截取
				if runtime.GOOS != "windows" {
					file = file[len(pwd()):]
				}

				return fmt.Sprintf("%s:%d", file, line)
			},
			timeFmt: func(layout string) string {
				now := time.Now()
				if layout == "" {
					layout = FmtYmdHmsSSS
				}
				return now.Format(layout)
			},
			levelFmt: func(i any) string {
				return strings.ToUpper(fmt.Sprintf(" [%s] ", i))
			},
		},
		Keys:    make(map[string]any),
		appName: hostName,
		buffer: sync.Pool{
			New: func() any {
				return &bytes.Buffer{}
			},
		},
	}
	l.writer = multiWriter{[]logWriter{&syncWriter{w: os.Stdout, l: l}}}
	l.formatter.caller = func(skip int) string {
		_, file, line, _ := runtime.Caller(skip)
		return fmt.Sprintf("%s:%d", file, line)
	}
	l.buffer.Put(bytes.NewBufferString(""))
	return l
}

func NewWithCtx(ctx context.Context) Logger {
	l := New()
	l.ctx = ctx
	return l
}

func (l *Logger) WithContext(ctx context.Context) Logger {
	l.ctx = ctx
	return *l
}
func (l Logger) AppName(appName string) Logger {
	l.appName = appName
	return l
}

func (l Logger) FormatCaller(fun CallerFun) Logger {
	l.formatter.caller = fun
	return l
}

func (l Logger) CallerSkip(skip int) Logger {
	l.callerSkip = skip
	return l
}

func (l Logger) LevelFormatter(fun FormatLevelFun) Logger {
	l.formatter.levelFmt = fun
	return l
}
func (l Logger) Output(w io.Writer) Logger {
	w1 := l.writer.writers
	w1 = append(w1, syncWriter{w, l})
	l.writer.writers = w1
	return l
}

func (l Logger) MultiWriter(writers ...io.Writer) Logger {
	var logWriters []logWriter
	for _, writer := range writers {
		logWriters = append(logWriters, syncWriter{writer, l})
	}
	l.writer.writers = logWriters
	return l
}

func (l Logger) Level(level Level) Logger {
	l.level = level
	return l
}

func (l *Logger) Dict(key string, value any) {
	c := l.ctx
	if c == nil {
		c = context.TODO()
	}
	c = context.WithValue(c, key, value)
	l.Keys[key] = value
	l.ctx = c
}

func (l Logger) dict() string {
	if len(l.Keys) > 0 {
		if data, e := json.Marshal(l.Keys); e != nil {
			return ""
		} else {
			return string(data)
		}
	}
	return ""
}

func (l Logger) traceId() string {
	c := l.ctx
	if c == nil {
		return ""
	}
	v := c.Value(TRACEID)
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%s", v)
}

func (l Logger) msg(lvl Level, msgFmt string, args ...any) string {
	caller := l.formatter.caller(l.callerSkip)
	time := l.formatter.timeFmt(l.FormatTime)
	level := l.formatter.levelFmt(lvl.GetName())
	msg := fmt.Sprintf(msgFmt, args...)
	msg = fmt.Sprintf("%s %s %s [%s] %s %s %s : %s\n", time, hostName, l.appName, l.traceId(), level, caller, l.dict(), msg)
	return msg
}

func (l Logger) WriteLevel(lvl Level, msgFmt string, args ...any) error {
	msg := l.msg(lvl, msgFmt, args...)
	b := l.buffer.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteString(msg)
	defer l.buffer.Put(b)
	_, err := l.writer.Write(b.Bytes())
	return err
}
