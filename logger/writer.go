package logger

import (
	"fmt"
	"io"
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
	writers []logWriter
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
