package logi

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/p9c/pod/pkg/util/logi/Pkg/Pk"
)

const (
	Off   = "off"
	Fatal = "fatal"
	Error = "error"
	Warn  = "warn"
	Info  = "info"
	Check = "check"
	Debug = "debug"
	Trace = "trace"
)

var (
	Levels = []string{
		Off,
		Fatal,
		Error,
		Check,
		Warn,
		Info,
		Debug,
		Trace,
	}
	Tags = map[string]string{
		Off:   "",
		Fatal: "FTL",
		Error: "ERR",
		Check: "CHK",
		Warn:  "WRN",
		Info:  "INF",
		Debug: "DBG",
		Trace: "TRC",
	}
	LevelsMap = map[string]int{
		Off:   0,
		Fatal: 1,
		Error: 2,
		Check: 3,
		Warn:  4,
		Info:  5,
		Debug: 6,
		Trace: 7,
	}
	StartupTime = time.Now()
)

type LogWriter struct {
	io.Writer
	write bool
}

// DirectionString is a helper function that returns a string that represents the direction of a connection (inbound or outbound).
func DirectionString(inbound bool) string {
	if inbound {
		return "inbound"
	}
	return "outbound"
}

func PickNoun(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

func (w *LogWriter) Print(a ...interface{}) {
	if w.write {
		_, _ = fmt.Fprint(w.Writer, a...)
	}
}

func (w *LogWriter) Printf(format string, a ...interface{}) {
	if w.write {
		_, _ = fmt.Fprintf(w.Writer, format, a...)
	}
}

func (w *LogWriter) Println(a ...interface{}) {
	if w.write {
		_, _ = fmt.Fprintln(w.Writer, a...)
	}
}

// Entry is a log entry to be printed as json to the log file
type Entry struct {
	Time         time.Time
	Level        string
	Package      string
	CodeLocation string
	Text         string
}

type (
	PrintcFunc  func(pkg string, fn func() string)
	PrintfFunc  func(pkg string, format string, a ...interface{})
	PrintlnFunc func(pkg string, a ...interface{})
	CheckFunc   func(pkg string, err error) bool
	SpewFunc    func(pkg string, a interface{})

	// Logger is a struct containing all the functions with nice handy names
	Logger struct {
		Packages      *Pk.Package
		Level         string
		Fatal         PrintlnFunc
		Error         PrintlnFunc
		Warn          PrintlnFunc
		Info          PrintlnFunc
		Check         CheckFunc
		Debug         PrintlnFunc
		Trace         PrintlnFunc
		Fatalf        PrintfFunc
		Errorf        PrintfFunc
		Warnf         PrintfFunc
		Infof         PrintfFunc
		Debugf        PrintfFunc
		Tracef        PrintfFunc
		Fatalc        PrintcFunc
		Errorc        PrintcFunc
		Warnc         PrintcFunc
		Infoc         PrintcFunc
		Debugc        PrintcFunc
		Tracec        PrintcFunc
		Fatals        SpewFunc
		Errors        SpewFunc
		Warns         SpewFunc
		Infos         SpewFunc
		Debugs        SpewFunc
		Traces        SpewFunc
		LogFileHandle *os.File
		Writer        LogWriter
		Color         bool
		Split         string
		LogChan       []chan Entry
	}
)

var L = NewLogger()

// AddLogChan adds a channel that log entries are sent to
func (l *Logger) AddLogChan() (ch chan Entry) {
	L.LogChan = append(L.LogChan, make(chan Entry))
	// L.Write = false
	return L.LogChan[len(L.LogChan)-1]
}

func NewLogger() (l *Logger) {
	p := make(Pk.Package)
	l = &Logger{
		Packages:      &p,
		Level:         "trace",
		LogFileHandle: os.Stderr,
		Color:         true,
		Split:         "pod",
		LogChan:       nil,
	}
	l.Fatal = l.printlnFunc(Fatal)
	l.Error = l.printlnFunc(Error)
	l.Warn = l.printlnFunc(Warn)
	l.Info = l.printlnFunc(Info)
	l.Check = l.checkFunc(Check)
	l.Debug = l.printlnFunc(Debug)
	l.Trace = l.printlnFunc(Trace)
	l.Fatalf = l.printfFunc(Fatal)
	l.Errorf = l.printfFunc(Error)
	l.Warnf = l.printfFunc(Warn)
	l.Infof = l.printfFunc(Info)
	l.Debugf = l.printfFunc(Debug)
	l.Tracef = l.printfFunc(Trace)
	l.Fatalc = l.printcFunc(Fatal)
	l.Errorc = l.printcFunc(Error)
	l.Warnc = l.printcFunc(Warn)
	l.Infoc = l.printcFunc(Info)
	l.Debugc = l.printcFunc(Debug)
	l.Tracec = l.printcFunc(Trace)
	l.Fatals = l.spewFunc(Fatal)
	l.Errors = l.spewFunc(Error)
	l.Warns = l.spewFunc(Warn)
	l.Infos = l.spewFunc(Info)
	l.Debugs = l.spewFunc(Debug)
	l.Traces = l.spewFunc(Trace)

	return
}

func (w *LogWriter) SetLogWriter(wr io.Writer) {
	w.Writer = wr
}

// SetLogPaths sets a file path to write logs
func (l *Logger) SetLogPaths(logPath, logFileName string) {
	const timeFormat = "2006-01-02_15-04-05"
	path := filepath.Join(logFileName, logPath)
	var logFileHandle *os.File
	if FileExists(path) {
		err := os.Rename(path, filepath.Join(logPath,
			time.Now().Format(timeFormat)+".json"))
		if err != nil {
			if L.Writer.write {
				L.Writer.Println("error rotating log", err)
			}
			return
		}
	}
	logFileHandle, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		if L.Writer.write {
			L.Writer.Println("error opening log file", logFileName)
		}
	}
	l.LogFileHandle = logFileHandle
	_, _ = fmt.Fprintln(logFileHandle, "{")
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func (l *Logger) SetLevel(level string, color bool, split string) {
	l.Level = sanitizeLoglevel(level)
	sep := string(os.PathSeparator)
	if runtime.GOOS == "windows" {
		sep = "\\"
		color = false
	}
	l.Split = split + sep
	l.Color = color
}

func (l *Logger) LocToPkg(pkg string) (out string) {
	// fmt.Println("pkg",pkg)
	sep := string(os.PathSeparator)
	if runtime.GOOS == "windows" {
		sep = "\\"
	}
	split := strings.Split(pkg, l.Split)
	if len(split) < 2 {
		return pkg
	}
	// fmt.Println("split",split, l.Split)
	pkg = split[1]
	split = strings.Split(pkg, sep)
	return strings.Join(split[:len(split)-1], string(os.PathSeparator))
}

func (l *Logger) Register(pkg string) string {
	// split := strings.Split(pkg, l.Split)
	// pkg = split[1]
	// split = strings.Split(pkg, string(os.PathSeparator))
	// pkg = strings.Join(split[:len(split)-1], string(os.PathSeparator))
	// if runtime.GOOS == "windows" {
	// 	pkg = strings.Replace(pkg, "/", string(os.PathSeparator), -1)
	// }
	pkg = l.LocToPkg(pkg)
	(*l.Packages)[pkg] = true
	return pkg
}

func (l *Logger) LoadConfig(configFile []byte) {
	var p Pk.Package
	if err := json.Unmarshal(configFile, &p); !l.Check("internal", err) {
		*l.Packages = p
	}
}

func init() {
	L.SetLevel("info", true, "pod")
	L.Writer.SetLogWriter(os.Stderr)
	L.Writer.write = true
	L.Trace("starting up logger")
}

func sanitizeLoglevel(level string) string {
	found := false
	for i := range Levels {
		if level == Levels[i] {
			found = true
			break
		}
	}
	if !found {
		level = "info"
	}
	return level
}

func trimReturn(s string) string {
	if s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}

func (l *Logger) LevelIsActive(level string) (out bool) {
	if LevelsMap[l.Level] >= LevelsMap[level] {
		out = true
	}
	return
}

func (l *Logger) GetLoc(loc string, line int) (out string) {
	split := strings.Split(loc, l.Split)
	if len(split) < 2 {
		out = loc
	} else {
		out = split[1]
	}
	return out + fmt.Sprint(":", line)
}

// printfFunc prints a log entry with formatting
func (l *Logger) printfFunc(level string) PrintfFunc {
	f := func(pkg, format string, a ...interface{}) {
		text := fmt.Sprintf(format, a...)
		if !l.LevelIsActive(level) || !(*l.Packages)[pkg] {
			return
		}
		if l.Writer.write || (*l.Packages)[pkg] {
			l.Writer.Println(Composite(text, level))
		}
		if l.LogChan != nil {
			_, loc, line, _ := runtime.Caller(2)
			pkg := l.LocToPkg(loc)
			out := Entry{time.Now(), level,
				pkg, l.GetLoc(loc, line), text}
			for i := range l.LogChan {
				l.LogChan[i] <- out
			}
		}
	}
	return f
}

// printcFunc prints from a closure returning a string
func (l *Logger) printcFunc(level string) PrintcFunc {
	f := func(pkg string, fn func() string) {
		if !l.LevelIsActive(level) || !(*l.Packages)[pkg] {
			return
		}
		t := fn()
		text := trimReturn(t)
		if l.Writer.write {
			l.Writer.Println(Composite(text, level))
		}
		if l.LogChan != nil {
			_, loc, line, _ := runtime.Caller(2)
			pkg := l.LocToPkg(loc)
			out := Entry{time.Now(), level,
				pkg, l.GetLoc(loc, line), text}
			for i := range l.LogChan {
				l.LogChan[i] <- out
			}
		}
	}
	return f
}

// printlnFunc prints a log entry like Println
func (l *Logger) printlnFunc(level string) PrintlnFunc {
	f := func(pkg string, a ...interface{}) {
		if !l.LevelIsActive(level) || !(*l.Packages)[pkg] {
			return
		}
		text := trimReturn(fmt.Sprintln(a...))
		if l.Writer.write {
			l.Writer.Println(Composite(text, l.Level))
		}
		if l.LogChan != nil {
			_, loc, line, _ := runtime.Caller(2)
			pkg := l.LocToPkg(loc)
			out := Entry{time.Now(), level, pkg,
				l.GetLoc(loc, line), text}
			for i := range l.LogChan {
				l.LogChan[i] <- out
			}
		}
	}
	return f
}

func (l *Logger) checkFunc(level string) CheckFunc {
	f := func(pkg string, err error) (out bool) {
		if !l.LevelIsActive(level) || !(*l.Packages)[pkg] {
			return
		}
		n := err == nil
		if n {
			return false
		}
		text := err.Error()
		if l.Writer.write {
			l.Writer.Println(Composite(text, "CHK"))
		}
		if l.LogChan != nil {
			_, loc, line, _ := runtime.Caller(2)
			pkg := l.LocToPkg(loc)
			out := Entry{time.Now(), level,
				pkg, l.GetLoc(loc, line), text}
			for i := range l.LogChan {
				l.LogChan[i] <- out
			}
		}
		return true
	}
	return f
}

// spewFunc spews a variable
func (l *Logger) spewFunc(level string) SpewFunc {
	f := func(pkg string, a interface{}) {
		if !l.LevelIsActive(level) || !(*l.Packages)[pkg] {
			return
		}
		text := trimReturn(spew.Sdump(a))
		o := "" + Composite("spew:", level)
		o += "\n" + text + "\n"
		if l.Writer.write {
			l.Writer.Print(o)
		}
		if l.LogChan != nil {
			_, loc, line, _ := runtime.Caller(2)
			pkg := l.LocToPkg(loc)
			out := Entry{time.Now(), level, pkg,
				l.GetLoc(loc, line), text}
			for i := range l.LogChan {
				l.LogChan[i] <- out
			}

		}
	}
	return f
}

func Composite(text, level string) (final string) {
	skip := 3
	if level == Check {
		skip = 4
	}
	_, loc, iLine, _ := runtime.Caller(skip)
	line := fmt.Sprint(iLine)
	since := fmt.Sprintf("%v", time.Now().Sub(StartupTime)/time.Millisecond*time.Millisecond)
	final = Tags[level] + " " + since + " " + text + " " + loc + ":" + line
	return final
}
