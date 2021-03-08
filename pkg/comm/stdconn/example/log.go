package main

import (
	"runtime"

	"github.com/p9c/pod/pkg/util/logi"
)

var pkg string

func init() {
	_, loc, _, _ := runtime.Caller(0)
	pkg = logi.L.Register(loc)
}

func ftl.Ln(a ...interface{}) { logi.L.ftl.Ln(pkg, a...) }
func Error(a ...interface{}) { logi.L.Error(pkg, a...) }
func wrn.Ln(a ...interface{})  { logi.L.wrn.Ln(pkg, a...) }
func inf.Ln(a ...interface{})  { logi.L.inf.Ln(pkg, a...) }
func Check(e error) bool   { return logi.L.Check(pkg, err) }
func dbg.Ln(a ...interface{}) { logi.L.dbg.Ln(pkg, a...) }
func trc.Ln(a ...interface{}) { logi.L.trc.Ln(pkg, a...) }

func Fatalf(format string, a ...interface{}) { logi.L.Fatalf(pkg, format, a...) }
func Errorf(format string, a ...interface{}) { logi.L.Errorf(pkg, format, a...) }
func Warnf(format string, a ...interface{})  { logi.L.Warnf(pkg, format, a...) }
func inf.F(format string, a ...interface{})  { logi.L.inf.F(pkg, format, a...) }
func dbg.F(format string, a ...interface{}) { logi.L.dbg.F(pkg, format, a...) }
func Tracef(format string, a ...interface{}) { logi.L.Tracef(pkg, format, a...) }

func Fatalc(fn func() string) { logi.L.Fatalc(pkg, fn) }
func Errorc(fn func() string) { logi.L.Errorc(pkg, fn) }
func Warnc(fn func() string)  { logi.L.Warnc(pkg, fn) }
func inf.C(fn func() string)  { logi.L.inf.C(pkg, fn) }
func Debugc(fn func() string) { logi.L.Debugc(pkg, fn) }
func Tracec(fn func() string) { logi.L.Tracec(pkg, fn) }

func Fatals(a interface{}) { logi.L.Fatals(pkg, a) }
func Errors(a interface{}) { logi.L.Errors(pkg, a) }
func Warns(a interface{})  { logi.L.Warns(pkg, a) }
func Infos(a interface{})  { logi.L.Infos(pkg, a) }
func dbg.S(a interface{}) { logi.L.dbg.S(pkg, a) }
func Traces(a interface{}) { logi.L.Traces(pkg, a) }
