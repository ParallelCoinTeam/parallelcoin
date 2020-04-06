package logi

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/pkg/logi/Pkg/Pk"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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
	StartupTime    = time.Now()
	BackgroundGrey = "\u001b[48;5;240m"
	ColorBlue      = "\u001b[38;5;33m"
	ColorBold      = "\u001b[1m"
	ColorBrown     = "\u001b[38;5;130m"
	ColorCyan      = "\u001b[36m"
	ColorFaint     = "\u001b[2m"
	ColorGreen     = "\u001b[38;5;40m"
	ColorItalic    = "\u001b[3m"
	ColorOff       = "\u001b[0m"
	ColorOrange    = "\u001b[38;5;208m"
	ColorPurple    = "\u001b[38;5;99m"
	ColorRed       = "\u001b[38;5;196m"
	ColorUnderline = "\u001b[4m"
	ColorViolet    = "\u001b[38;5;201m"
	ColorYellow    = "\u001b[38;5;226m"
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
	//L.Write = false
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

func (wr *LogWriter) SetLogWriter(w io.Writer) {
	wr.Writer = w
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
		sep = "/"
	}
	l.Split = split + sep
	l.Color = color
}

func (l *Logger) LocToPkg(pkg string) (out string) {
	//fmt.Println("pkg",pkg)
	sep := string(os.PathSeparator)
	if runtime.GOOS == "windows" {
		sep = "/"
	}
	split := strings.Split(pkg, l.Split)
	//fmt.Println("split",split, l.Split)
	pkg = split[1]
	split = strings.Split(pkg, sep)
	return strings.Join(split[:len(split)-1], string(os.PathSeparator))
}

func (l *Logger) Register(pkg string) string {
	//split := strings.Split(pkg, l.Split)
	//pkg = split[1]
	//split = strings.Split(pkg, string(os.PathSeparator))
	//pkg = strings.Join(split[:len(split)-1], string(os.PathSeparator))
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
	L.SetLevel("trace", true, "pod")
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

var TermWidth = func() int { return 80 }

func (l *Logger) GetLoc(loc string, line int) (out string) {
	split := strings.Split(loc, l.Split)
	if len(split) < 2 {
		out = loc
	}
	out = split[1]
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
			l.Writer.Println(Composite(text, level, l.Color, l.Split))
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
			l.Writer.Println(Composite(text, level, l.Color, l.Split))
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
			l.Writer.Println(Composite(text, l.Level, l.Color, l.Split))
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
			l.Writer.Println(Composite(text, "CHK", l.Color, l.Split))
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
		o := "" + Composite("spew:", level, l.Color, l.Split)
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

func Composite(text, level string, color bool, split string) string {
	dots := "."
	terminalWidth := TermWidth()
	if TermWidth() <= 120 {
		terminalWidth = 120
	}
	skip := 3
	if level == Check {
		skip = 4
	}
	_, loc, iline, _ := runtime.Caller(skip)
	line := fmt.Sprint(iline)
	files := strings.Split(loc, split)
	var file, since string
	file = loc
	if len(files) > 1 {
		file = files[1]
	}
	switch {
	case terminalWidth <= 60:
		since = ""
		file = ""
		line = ""
		dots = " "
	case terminalWidth <= 80:
		dots = " "
		if len(file) > 30 {
			file = ""
			line = ""
		}
		since = fmt.Sprintf("%v", time.Now().Sub(StartupTime)/time.Second*time.Second)
	case terminalWidth < 120:
		if len(file) > 40 {
			file = ""
			line = ""
			dots = " "
		}
		since = fmt.Sprintf("%v", time.Now().Sub(StartupTime)/time.Millisecond*time.Millisecond)
	case terminalWidth < 160:
		if len(file) > 60 {
			file = ""
			line = ""
			dots = " "
		}
		since = fmt.Sprintf("%v", time.Now().Sub(StartupTime)/time.Millisecond*time.Millisecond)
		//since = fmt.Sprint(time.Now())[:19]
	case terminalWidth >= 200:
		since = fmt.Sprint(time.Now())[:39]
	default:
		since = fmt.Sprint(time.Now())[:19]
	}
	levelLen := 4 // len(level) + 1
	sinceLen := len(since) + 1
	textLen := len(text) + 1
	fileLen := len(file) + 1
	lineLen := len(line) + 1
	if file != "" {
		file += ":"
	}
	if color {
		switch Tags[level] {
		case "FTL":
			level = ColorBold + ColorRed + Tags[level] + ColorOff
			since = ColorRed + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "ERR":
			level = ColorBold + ColorOrange + Tags[level] + ColorOff
			since = ColorOrange + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "WRN":
			level = ColorBold + ColorYellow + Tags[level] + ColorOff
			since = ColorYellow + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "INF":
			level = ColorBold + ColorGreen + Tags[level] + ColorOff
			since = ColorGreen + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "CHK":
			level = ColorBold + ColorCyan + Tags[level] + ColorOff
			since = since
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "DBG":
			level = ColorBold + ColorBlue + Tags[level] + ColorOff
			since = ColorBlue + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "TRC":
			level = ColorBold + ColorViolet + Tags[level] + ColorOff
			since = ColorViolet + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		}
	} else {
		level = Tags[level]
	}
	final := ""
	if levelLen+sinceLen+textLen+fileLen+lineLen > terminalWidth {
		lines := strings.Split(text, "\n")
		// log text is multiline
		line1len := terminalWidth - levelLen - sinceLen - fileLen - lineLen
		restLen := terminalWidth - levelLen - sinceLen
		if len(lines) > 1 {
			final = fmt.Sprintf("%s %s %s %s%s", level, since,
				strings.Repeat(".",
					terminalWidth-levelLen-sinceLen-fileLen-lineLen-4),
				file, line)
			final += text[:len(text)-1]
		} else {
			// log text is a long line
			spaced := strings.Split(text, " ")
			var rest bool
			curLineLen := 0
			final += fmt.Sprintf("%s %s ", level, since)
			var i int
			for i = range spaced {
				if i > 0 {
					curLineLen += len(spaced[i-1]) + 1
					if !rest {
						if curLineLen >= line1len {
							rest = true
							spacers := terminalWidth - levelLen - sinceLen -
								fileLen - lineLen - curLineLen + len(spaced[i-1]) + 1
							if spacers < 1 {
								spacers = 1
							}
							final += strings.Repeat(dots, spacers)
							final += fmt.Sprintf(" %s%s\n",
								file, line)
							final += strings.Repeat(" ", levelLen+sinceLen)
							final += spaced[i-1] + " "
							curLineLen = len(spaced[i-1]) + 1
						} else {
							final += spaced[i-1] + " "
						}
					} else {
						if curLineLen >= restLen-1 {
							final += "\n" + strings.Repeat(" ",
								levelLen+sinceLen)
							final += spaced[i-1] + dots
							curLineLen = len(spaced[i-1]) + 1
						} else {
							final += spaced[i-1] + " "
						}
					}
				}
			}
			curLineLen += len(spaced[i])
			if !rest {
				if curLineLen >= line1len {
					final += fmt.Sprintf("%s %s%s\n",
						strings.Repeat(dots,
							len(spaced[i])+line1len-curLineLen),
						file, line)
					final += strings.Repeat(" ", levelLen+sinceLen)
					final += spaced[i] // + "\n"
				} else {
					final += fmt.Sprintf("%s %s %s%s\n",
						spaced[i],
						strings.Repeat(dots,
							terminalWidth-curLineLen-fileLen-lineLen),
						file, line)
				}
			} else {
				if curLineLen >= restLen {
					final += "\n" + strings.Repeat(" ", levelLen+sinceLen)
				}
				final += spaced[i]
			}
		}
	} else {
		final = fmt.Sprintf("%s %s %s %s %s%s", level, since, text,
			strings.Repeat(dots,
				terminalWidth-levelLen-sinceLen-textLen-fileLen-lineLen),
			file, line)
	}
	return final
}
