package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type logWriter interface {
	io.Writer
	LevelMsg(lvl Level, msg string) error
}

type syncWriter struct {
	w io.Writer
	l Logger
}

func (w syncWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w syncWriter) LevelMsg(lvl Level, msg string) error {
	caller := w.l.formatter.caller(w.l.callerSkip)
	time := w.l.formatter.timeFmt(w.l.FormatTime)
	level := w.l.formatter.levelFmt(lvl.GetName())
	msg = fmt.Sprintf("%s %s %s [%s] %s %s : %s", time, hostName, w.l.appName, w.l.traceId(), level, caller, msg)
	return w.write(msg)
}

func (w syncWriter) write(msg string) error {
	_, e := w.Write([]byte(msg))
	return e
}

func (w syncWriter) Msg(msg string) error {
	return w.write(msg)
}

type multiWriter struct {
	writers []io.Writer
}

func (m multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range m.writers {
		if _n, _err := w.Write(p); err == nil {
			n = _n
			if _err != nil {
				err = _err
			} else if _n != len(p) {
				err = io.ErrShortWrite
			}
		}
	}
	return n, err
}

type FileLevelWriter struct {
	*os.File
	link     string
	original string
}

//NewFileLevelWriter return file writer,it's name is appName-lvl-time.log,eg: myApp-info.log and link myApp-info-20230208.${idx}.log
func (l Logger) NewFileLevelWriter(lvl Level) *FileLevelWriter {
	w := &FileLevelWriter{}
	linkName := l.appName + "-" + lvl.GetName() + ".log"
	original := l.appName + "-" + lvl.GetName() + time.Now().Format(timeLayout) + ".log"
	f := func(dst string, log Logger) *os.File {
		f, e := os.Create(dst)
		if e != nil {
			log.Error("无法创建日志文件,%s,%v", dst, e)
			return nil
		} else {
			return f
		}
	}
	writer := f(original, l)
	w.File = writer
	w.original = original
	w.link = linkName
	_ = os.Symlink(original, linkName)
	return w
}
