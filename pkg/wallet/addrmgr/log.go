package waddrmgr

import (
	"github.com/parallelcointeam/parallelcoin/pkg/log"
)

type _dtype int

var _d _dtype
var l = log.NewLogger("info")

func UseLogger(logger *log.Logger) {
	l = logger
}

var (
	FATAL  = l.Fatal
	ERROR  = l.Error
	WARN   = l.Warn
	INFO   = l.Info
	DEBUG  = l.Debug
	TRACE  = l.Trace
	FATALF = l.Fatalf
	ERRORF = l.Errorf
	WARNF  = l.Warnf
	INFOF  = l.Infof
	DEBUGF = l.Debugf
	TRACEF = l.Tracef
	FATALC = l.Fatalc
	ERRORC = l.Errorc
	WARNC  = l.Warnc
	INFOC  = l.Infoc
	DEBUGC = l.Debugc
	TRACEC = l.Tracec
)
