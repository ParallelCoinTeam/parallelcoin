package main

import (
	"github.com/p9c/pod/pkg/util/logi"
	"runtime"
)

var pkg string

func init() {
	_, loc, _, _ := runtime.Caller(0)
	pkg = logi.L.Register(loc)
}

func Fatal(a ...interface{}) { logi.L.Fatal(pkg, a...) }
func Error(a ...interface{}) { logi.L.Error(pkg, a...) }
func Warn(a ...interface{})  { logi.L.Warn(pkg, a...) }
func Info(a ...interface{})  { logi.L.Info(pkg, a...) }
func Check(err error) bool   { return logi.L.Check(pkg, err) }
func Debug(a ...interface{}) { logi.L.Debug(pkg, a...) }
func Trace(a ...interface{}) { logi.L.Trace(pkg, a...) }

func Fatalf(format string, a ...interface{}) { logi.L.Fatalf(pkg, format, a...) }
func Errorf(format string, a ...interface{}) { logi.L.Errorf(pkg, format, a...) }
func Warnf(format string, a ...interface{})  { logi.L.Warnf(pkg, format, a...) }
func Infof(format string, a ...interface{})  { logi.L.Infof(pkg, format, a...) }
func Debugf(format string, a ...interface{}) { logi.L.Debugf(pkg, format, a...) }
func Tracef(format string, a ...interface{}) { logi.L.Tracef(pkg, format, a...) }

func Fatalc(fn func() string) { logi.L.Fatalc(pkg, fn) }
func Errorc(fn func() string) { logi.L.Errorc(pkg, fn) }
func Warnc(fn func() string)  { logi.L.Warnc(pkg, fn) }
func Infoc(fn func() string)  { logi.L.Infoc(pkg, fn) }
func Debugc(fn func() string) { logi.L.Debugc(pkg, fn) }
func Tracec(fn func() string) { logi.L.Tracec(pkg, fn) }

func Fatals(a interface{}) { logi.L.Fatals(pkg, a) }
func Errors(a interface{}) { logi.L.Errors(pkg, a) }
func Warns(a interface{})  { logi.L.Warns(pkg, a) }
func Infos(a interface{})  { logi.L.Infos(pkg, a) }
func Debugs(a interface{}) { logi.L.Debugs(pkg, a) }
func Traces(a interface{}) { logi.L.Traces(pkg, a) }

var commands = map[string][]string{
	"build": {
		"go build -v",
	},
	"windows": {
		`go build -v -ldflags="-H windowsgui \"%ldflags"\"`,
	},
	"tests": {
		"go test ./...",
	},
	"kopachgui": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet -l debug --lan --solo --kopachgui kopach",
	},
	"gui": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet --lan",
	},
	"guis": {
		"go install -v %ldflags",
		"pod -D test1 --minerpass pa55word",
	},
	"guass": {
		"go install -v %ldflags",
		"pod -D %datadir -l trace -g 1 -G --solo --lan --minerpass pa55word",
	},
	"guihttpprof": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet --lan --solo --kopachgui --profile 6969",
	},
	"guiprof": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet --lan --solo --kopachgui",
	},
	"mainnode": {
		"go install -v %ldflags",
		"pod -D testmain -n mainnet -l info --connect seed3.parallelcoin." +
			"io:11047 node",
	},
	"mainwallet": {
		"go install -v %ldflags",
		"pod -D testmain -n mainnet -l trace wallet",
	},
	"teststopkopach": {
		"go install -v %ldflags",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D %datadir --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan kopach",
	},
	"teststopnode": {
		"go install -v %ldflags",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D %datadir --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan node",
	},
	"teststopwallet": {
		"go install -v %ldflags",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D %datadir --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan wallet",
	},
	"nodegui": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet nodegui",
	},
	"testnode": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet -l trace --solo --lan --norpc=false node",
	},
	"testwallet": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet -l trace --walletpass aoeuaoeu --solo --lan wallet",
	},
	"testkopach": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet -l trace -g -G 1 --solo --lan kopach",
	},
	"resetwallet": {
		"pod -D %datadir -l trace --walletpass aoeuaoeu wallet drophistory",
	},
	"stroy": {
		"go install -v %ldflags ./stroy/.",
	},
}

