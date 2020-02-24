package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	
	gt "github.com/buger/goterm"
	"github.com/davecgh/go-spew/spew"
)

var L = Empty()

var (
	colorRed    = "\u001b[38;5;196m"
	colorOrange = "\u001b[38;5;208m"
	colorYellow = "\u001b[38;5;226m"
	colorGreen  = "\u001b[38;5;40m"
	colorBlue   = "\u001b[38;5;33m"
	// colorPurple    = "\u001b[38;5;99m"
	colorViolet = "\u001b[38;5;201m"
	// colorBrown     = "\u001b[38;5;130m"
	colorBold = "\u001b[1m"
	// colorUnderline = "\u001b[4m"
	colorItalic = "\u001b[3m"
	// colorFaint     = "\u001b[2m"
	colorOff       = "\u001b[0m"
	backgroundGrey = "\u001b[48;5;240m"
)

var StartupTime = time.Now()

type PrintlnFunc *func(a ...interface{})
type PrintfFunc *func(format string, a ...interface{})
type PrintcFunc *func(func() string)
type SpewFunc *func(interface{})

const (
	Off   = "off"
	Fatal = "fatal"
	Error = "error"
	Warn  = "warn"
	Info  = "info"
	Debug = "debug"
	Trace = "trace"
)

var Levels = []string{
	Off, Fatal, Error, Warn, Info, Debug, Trace,
}

type LogWriter struct {
	io.Writer
}

var wr LogWriter

func init() {
	SetLogWriter(os.Stderr)
	L.SetLevel("info", true)
	TRACE("starting up logger")
}

func Print(a ...interface{}) {
	wr.Print(a...)
}

func Println(a ...interface{}) {
	wr.Println(a...)
}

func Printf(format string, a ...interface{}) {
	wr.Printf(format, a...)
}

// Logger is a struct containing all the functions with nice handy names
type Logger struct {
	Fatal         PrintlnFunc
	Error         PrintlnFunc
	Warn          PrintlnFunc
	Info          PrintlnFunc
	Debug         PrintlnFunc
	Trace         PrintlnFunc
	Traces        SpewFunc
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
	LogFileHandle *os.File
	Writer        LogWriter
	Color         bool
	// If this channel is loaded log entries are composed and sent to it
	LogChan chan Entry
}

// Entry is a log entry to be printed as json to the log file
type Entry struct {
	Time         time.Time
	Level        string
	CodeLocation string
	Text         string
}

func Empty() *Logger {
	return &Logger{
		Fatal:  NoPrintln(),
		Error:  NoPrintln(),
		Warn:   NoPrintln(),
		Info:   NoPrintln(),
		Debug:  NoPrintln(),
		Trace:  NoPrintln(),
		Traces: NoSpew(),
		Fatalf: NoPrintf(),
		Errorf: NoPrintf(),
		Warnf:  NoPrintf(),
		Infof:  NoPrintf(),
		Debugf: NoPrintf(),
		Tracef: NoPrintf(),
		Fatalc: NoClosure(),
		Errorc: NoClosure(),
		Warnc:  NoClosure(),
		Infoc:  NoClosure(),
		Debugc: NoClosure(),
		Tracec: NoClosure(),
		Writer: wr,
	}
	
}

// sanitizeLoglevel accepts a string and returns a
// default if the input is not in the Levels slice
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

func SetLogWriter(w io.Writer) {
	wr.Writer = w
}

func (w *LogWriter) Print(a ...interface{}) {
	_, _ = fmt.Fprint(wr, a...)
}

func (w *LogWriter) Println(a ...interface{}) {
	_, _ = fmt.Fprintln(wr, a...)
}

func (w *LogWriter) Printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(wr, format, a...)
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
			wr.Println("error rotating log", err)
			return
		}
	}
	logFileHandle, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		wr.Println("error opening log file", logFileName)
	}
	l.LogFileHandle = logFileHandle
	_, _ = fmt.Fprintln(logFileHandle, "{")
}

// SetLevel enables or disables the various print functions
func (l *Logger) SetLevel(level string, color bool) *Logger {
	// *l = *Empty()
	level = sanitizeLoglevel(level)
	var fallen bool
	switch {
	case level == Trace || fallen:
		TRACE("trace testing")
		l.Trace = printlnFunc("TRC", color, l.LogFileHandle, nil)
		l.Tracef = printfFunc("TRC", color, l.LogFileHandle, nil)
		l.Tracec = printcFunc("TRC", color, l.LogFileHandle, nil)
		l.Traces = ps("TRC", color, l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Debug || fallen:
		l.Debug = printlnFunc("DBG", color, l.LogFileHandle, l.LogChan)
		l.Debugf = printfFunc("DBG", color, l.LogFileHandle, l.LogChan)
		l.Debugc = printcFunc("DBG", color, l.LogFileHandle, l.LogChan)
		fallen = true
		fallthrough
	case level == Info || fallen:
		l.Info = printlnFunc("INF", color, l.LogFileHandle, l.LogChan)
		l.Infof = printfFunc("INF", color, l.LogFileHandle, l.LogChan)
		l.Infoc = printcFunc("INF", color, l.LogFileHandle, l.LogChan)
		fallen = true
		fallthrough
	case level == Warn || fallen:
		l.Warn = printlnFunc("WRN", color, l.LogFileHandle, l.LogChan)
		l.Warnf = printfFunc("WRN", color, l.LogFileHandle, l.LogChan)
		l.Warnc = printcFunc("WRN", color, l.LogFileHandle, l.LogChan)
		fallen = true
		fallthrough
	case level == Error || fallen:
		l.Error = printlnFunc("ERR", color, l.LogFileHandle, l.LogChan)
		l.Errorf = printfFunc("ERR", color, l.LogFileHandle, l.LogChan)
		l.Errorc = printcFunc("ERR", color, l.LogFileHandle, l.LogChan)
		fallen = true
		fallthrough
	case level == Fatal:
		l.Fatal = printlnFunc("FTL", color, l.LogFileHandle, l.LogChan)
		l.Fatalf = printfFunc("FTL", color, l.LogFileHandle, l.LogChan)
		l.Fatalc = printcFunc("FTL", color, l.LogFileHandle, l.LogChan)
		fallen = true
	}
	return l
}

var NoPrintln = func() PrintlnFunc {
	f := func(_ ...interface{}) {}
	return &f
}
var NoPrintf = func() PrintfFunc {
	f := func(_ string, _ ...interface{}) {}
	return &f
}
var NoClosure = func() PrintcFunc {
	f := func(_ func() string) {}
	return &f
}
var NoSpew = func() SpewFunc {
	f := func(_ interface{}) {}
	return &f
}

func trimReturn(s string) string {
	if s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}

func Composite(text, level string, color bool) string {
	terminalWidth := gt.Width()
	if terminalWidth < 120 {
		terminalWidth = 120
	}
	skip := 3
	if level == "ERR" {
		skip = 4
	}
	_, loc, iline, _ := runtime.Caller(skip)
	line := fmt.Sprint(iline)
	files := strings.Split(loc, "pod/")
	var file string
	if len(files) > 1 {
		file = files[1]
	}
	sinceW := 12
	if level == "ERR" {
		sinceW = 9
	}
	since := fmt.Sprintf("%-"+fmt.Sprint(sinceW)+"s", time.Now().Sub(StartupTime)/time.Second*time.Second)
	if terminalWidth > 200 {
		since = fmt.Sprint(time.Now())[:25]
	}
	levelLen := len(level) + 1
	sinceLen := len(since) + 1
	textLen := len(text) + 1
	fileLen := len(file) + 1
	lineLen := len(line) + 1
	if color {
		switch level {
		case "FTL":
			level = colorBold + colorRed + level + colorOff
			since = colorRed + since + colorOff
			file = colorItalic + colorBlue + file
			line = line + colorOff
		case "ERR":
			level = colorBold + colorOrange + level + colorOff
			since = colorOrange + since + colorOff
			file = colorItalic + colorBlue + file
			line = line + colorOff
		case "WRN":
			level = colorBold + colorYellow + level + colorOff
			since = colorYellow + since + colorOff
			file = colorItalic + colorBlue + file
			line = line + colorOff
		case "INF":
			level = colorBold + colorGreen + level + colorOff
			since = colorGreen + since + colorOff
			file = colorItalic + colorBlue + file
			line = line + colorOff
		case "DBG":
			level = colorBold + colorBlue + level + colorOff
			since = colorBlue + since + colorOff
			file = colorItalic + colorBlue + file
			line = line + colorOff
		case "TRC":
			level = colorBold + colorViolet + level + colorOff
			since = colorViolet + since + colorOff
			file = colorItalic + colorBlue + file
			line = line + colorOff
		case "CHK":
			level = colorBold + level + colorOff
			// since = since
			file = colorItalic + file
			line = line + colorOff
		}
	}
	final := ""
	if levelLen+sinceLen+textLen+fileLen+lineLen > terminalWidth {
		lines := strings.Split(text, "\n")
		// log text is multiline
		line1len := terminalWidth - levelLen - sinceLen - fileLen - lineLen
		restLen := terminalWidth - levelLen - sinceLen
		if len(lines) > 1 {
			final = fmt.Sprintf("%s %s %s %s:%s", level, since,
				strings.Repeat(" ",
					terminalWidth-levelLen-sinceLen-fileLen-lineLen),
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
							final += strings.Repeat(".", spacers)
							final += fmt.Sprintf(" %s:%s\n",
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
							final += spaced[i-1] + "."
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
					final += fmt.Sprintf("%s %s:%s\n",
						strings.Repeat(".",
							len(spaced[i])+line1len-curLineLen),
						file, line)
					final += strings.Repeat(" ", levelLen+sinceLen)
					final += spaced[i] // + "\n"
				} else {
					final += fmt.Sprintf("%s %s %s:%s\n",
						spaced[i],
						strings.Repeat(".",
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
		final = fmt.Sprintf("%s %s %s %s %s:%s", level, since, text,
			strings.Repeat(".",
				terminalWidth-levelLen-sinceLen-textLen-fileLen-lineLen),
			file, line)
	}
	return final
}

// printlnFunc prints a log entry like Println
func printlnFunc(level string, color bool, fh *os.File, ch chan Entry) PrintlnFunc {
	f := func(a ...interface{}) {
		text := trimReturn(fmt.Sprintln(a...))
		wr.Println(Composite(text, level, color))
		if fh != nil || ch != nil {
			_, loc, line, _ := runtime.Caller(2)
			out := Entry{time.Now(), level, fmt.Sprint(loc, ":", line), text}
			if fh != nil {
				j, err := json.Marshal(out)
				if err != nil {
					wr.Println("logging error:", err)
				}
				_, _ = fmt.Fprint(fh, string(j)+",")
			}
			if ch != nil {
				ch <- out
			}
		}
	}
	return &f
}

// printfFunc prints a log entry with formatting
func printfFunc(level string, color bool, fh *os.File, ch chan Entry) PrintfFunc {
	f := func(format string, a ...interface{}) {
		text := fmt.Sprintf(format, a...)
		wr.Println(Composite(text, level, color))
		if fh != nil || ch != nil {
			_, loc, line, _ := runtime.Caller(2)
			out := Entry{time.Now(), level, fmt.Sprint(loc, ":", line), text}
			if fh != nil {
				j, err := json.Marshal(out)
				if err != nil {
					wr.Println("logging error:", err)
				}
				_, _ = fmt.Fprint(fh, string(j)+",")
			}
			if ch != nil {
				ch <- out
			}
		}
	}
	return &f
}

// printcFunc prints from a closure returning a string
func printcFunc(level string, color bool, fh *os.File, ch chan Entry) PrintcFunc {
	f := func(fn func() string) {
		t := fn()
		text := trimReturn(t)
		wr.Println(Composite(text, level, color))
		if fh != nil || ch != nil {
			_, loc, line, _ := runtime.Caller(2)
			out := Entry{time.Now(), level, fmt.Sprint(loc, ":", line), text}
			if fh != nil {
				j, err := json.Marshal(out)
				if err != nil {
					wr.Println("logging error:", err)
				}
				_, _ = fmt.Fprint(fh, string(j)+",")
			}
			if ch != nil {
				ch <- out
			}
		}
	}
	return &f
}

// ps spews a variable
func ps(level string, color bool, fh *os.File) SpewFunc {
	f := func(a interface{}) {
		text := trimReturn(spew.Sdump(a))
		o := "" + Composite("spew:", level, color)
		o += "\n" + text + "\n"
		wr.Print(o)
		if fh != nil {
			_, loc, line, _ := runtime.Caller(2)
			out := Entry{time.Now(), level, fmt.Sprint(loc, ":", line), text}
			j, err := json.Marshal(out)
			if err != nil {
				wr.Println("logging error:", err)
			}
			_, _ = fmt.Fprint(fh, string(j)+",")
		}
	}
	return &f
}

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// DirectionString is a helper function that returns a string that represents the direction of a connection (inbound or outbound).
func DirectionString(inbound bool) string {
	if inbound {
		return "inbound"
	}
	return "outbound"
}

// PickNoun returns the singular or plural form of a noun depending
// on the count n.
func PickNoun(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
