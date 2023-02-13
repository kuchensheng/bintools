package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
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
	level    Level
	lock     sync.Mutex
}

//NewFileLevelWriter return file writer,it's name is appName-lvl-time.log,eg: myApp-info.log and link myApp-info-20230208.${idx}.log
func (l Logger) NewFileLevelWriter(lvl Level) *FileLevelWriter {
	w := &FileLevelWriter{}
	w.level = lvl
	if _, e := os.ReadDir(l.logHome); e != nil && os.IsNotExist(e) {
		os.MkdirAll(l.logHome, 666)
	}
	linkName := filepath.Join(l.logHome, l.appName+"-"+lvl.GetName()+".log")
	original := filepath.Join(l.logHome, l.appName+"-"+lvl.GetName()+time.Now().Format(timeLayout)+".log")
	w.original = original
	w.link = linkName
	w.lock = sync.Mutex{}
	return w
}

func (w *FileLevelWriter) CreateWriter(dst, link string, log Logger) {
	f, e := os.OpenFile(dst, os.O_CREATE|os.O_APPEND|os.O_RDWR, 666)
	if e != nil {
		log.Error("无法创建日志文件,%s,%v", dst, e)
	} else {
		_ = os.Symlink(dst, link)
		w.File = f
	}
}
