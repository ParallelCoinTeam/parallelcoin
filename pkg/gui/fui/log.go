package fui

import (
	"runtime"

	"github.com/p9c/pod/pkg/util/logi"
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
