package main

import (
	"github.com/parallelcointeam/parallelcoin/pkg/log"
)

type _dtype int

var _d _dtype
var l *log.Logger

func UseLogger(logger *log.Logger) {
	l = logger
}

var (
	FATAL  log.PrintlnFunc // = l.Fatal
	ERROR  log.PrintlnFunc // = l.Error
	WARN   log.PrintlnFunc // = l.Warn
	INFO   log.PrintlnFunc // = l.Info
	DEBUG  log.PrintlnFunc // = l.Debug
	TRACE  log.PrintlnFunc // = l.Trace
	FATALF log.PrintfFunc  // = l.Fatalf
	ERRORF log.PrintfFunc  // = l.Errorf
	WARNF  log.PrintfFunc  // = l.Warnf
	INFOF  log.PrintfFunc  // = l.Infof
	DEBUGF log.PrintfFunc  // = l.Debugf
	TRACEF log.PrintfFunc  // = l.Tracef
	FATALC log.Closure     // = l.Fatalc
	ERRORC log.Closure     // = l.Errorc
	WARNC  log.Closure     // = l.Warnc
	INFOC  log.Closure     // = l.Infoc
	DEBUGC log.Closure     // = l.Debugc
	TRACEC log.Closure     // = l.Tracec
)

// pickNoun returns the singular or plural form of a noun depending on the count n.
func pickNoun(n int, singular, plural string) string {

	if n == 1 {
		return singular
	}
	return plural
}

// directionString is a helper function that returns a string that represents the direction of a connection (inbound or outbound).
func directionString(inbound bool) string {
	if inbound {
		return "inbound"
	}
	return "outbound"
}
