//// +build ignore

package main

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	
	"github.com/p9c/pod/pkg/util/logi"
)

type command struct {
	name string
	args []string
}

var commands = map[string][]string{
	"build": {
		"go build -v",
	},
	"windows": {
		`go build -v -ldflags="-H windowsgui"`,
	},
	"tests": {
		`go test ./...`,
	},
	"kopachgui": {
		"go install -v",
		"pod -D test0 -n testnet -l debug --lan --solo --kopachgui kopach",
	},
	"testkopach": {
		"go install -v",
		"pod -D test0 -n testnet -l trace -g -G 1 --lan kopach",
	},
	"testnode": {
		"go install -v",
		"pod -D test0 -n testnet -l debug --solo --lan node",
	},
	"nodegui": {
		"go install -v",
		"pod -D test0 -n testnet nodegui",
	},
	"gui": {
		"go install -v",
		"pod -D test0 -n testnet --lan",
	},
	"guis": {
		"go install -v",
		"pod -D test1 --minerpass pa55word",
	},
	"guass": {
		"go install -v",
		"pod -D test0 --minerpass pa55word",
	},
	"resetwallet0": {
		"pod -D test0 -l trace --walletpass aoeuaoeu wallet drophistory",
	},
	"resetwallet1": {
		"pod -D test1 -l trace --walletpass aoeuaoeu wallet drophistory",
	},
	"guihttpprof": {
		"go install -v",
		"pod -D test0 -n testnet --lan --solo --kopachgui --profile 6969",
	},
	"guiprof": {
		"go install -v",
		"pod -D test0 -n testnet --lan --solo --kopachgui",
	},
	"mainnode": {
		"go install -v",
		"pod -D testmain -n mainnet -l info --connect seed3.parallelcoin." +
			"io:11047 node",
	},
	"testwallet": {
		"go install -v",
		"pod -D test0 -n testnet -l trace --walletpass aoeuaoeu wallet",
	},
	"mainwallet": {
		"go install -v",
		"pod -D testmain -n mainnet -l trace wallet",
	},
	"teststopkopach": {
		"go install -v",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D test0 --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan kopach",
	},
	"teststopnode": {
		"go install -v",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D test0 --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan node",
	},
	"teststopwallet": {
		"go install -v",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D test0 --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan wallet",
	},
}

func main() {
	if len(os.Args) > 1 {
		if list, ok := commands[os.Args[1]]; ok {
			for i := range list {
				Info("executing item", i, "of list", os.Args[1], list[i])
				split := strings.Split(list[i], " ")
				cmd := exec.Command(split[0], split[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); Check(err) {
				}
			}
		}
	} else {
		Error("no command requested, available:")
		
	}
}

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
