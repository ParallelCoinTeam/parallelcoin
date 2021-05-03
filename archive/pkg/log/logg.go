package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gookit/color"
	uberatomic "go.uber.org/atomic"
)

const (
	_Off = iota
	_Fatal
	_Error
	_Chek
	_Warn
	_Info
	_Debug
	_Trace
)

type (
	// LevelPrinter defines a set of terminal printing primitives that output with
	// extra data, time, log logLevelList, and code location
	LevelPrinter struct {
		// Ln prints lists of interfaces with spaces in between
		Ln func(a ...interface{})
		// F prints like fmt.Println surrounded by log details
		F func(format string, a ...interface{})
		// S prints a spew.Sdump for an interface slice
		S func(a ...interface{})
		// C accepts a function so that the extra computation can be avoided if it is
		// not being viewed
		C func(closure func() string)
		// Chk is a shortcut for printing if there is an error, or returning true
		Chk func(e error) bool
	}
	logLevelList struct {
		Off, Fatal, Error, Check, Warn, Info, Debug, Trace int32
	}
	LevelSpec struct {
		ID        int32
		Name      string
		Colorizer func(format string, a ...interface{}) string
	}

	// Entry is a log entry to be printed as json to the log file
	Entry struct {
		Time         time.Time
		Level        string
		Package      string
		CodeLocation string
		Text         string
	}
)

var (
	logger_started = time.Now()
	App            = "   pod"
	AppColorizer   = color.White.Sprint
	// sep is just a convenient shortcut for this very longwinded expression
	sep          = string(os.PathSeparator)
	currentLevel = uberatomic.NewInt32(logLevels.Info)
	// writer can be swapped out for any io.*writer* that you want to use instead of
	// stdout.
	writer io.Writer = os.Stderr
	// allSubsystems stores all of the package subsystem names found in the current
	// application
	allSubsystems []string
	// highlighted is a text that helps visually distinguish a log entry by category
	highlighted = make(map[string]struct{})
	// logFilter specifies a set of packages that will not pr logs
	logFilter = make(map[string]struct{})
	// mutexes to prevent concurrent map accesses
	highlightMx, _logFilterMx sync.Mutex
	// logLevels is a shorthand access that minimises possible Name collisions in the
	// dot import
	logLevels = logLevelList{
		Off:   _Off,
		Fatal: _Fatal,
		Error: _Error,
		Check: _Chek,
		Warn:  _Warn,
		Info:  _Info,
		Debug: _Debug,
		Trace: _Trace,
	}
	// LevelSpecs specifies the id, string name and color-printing function
	LevelSpecs = []LevelSpec{
		{logLevels.Off, "off  ", color.Bit24(0, 0, 0, false).Sprintf},
		{logLevels.Fatal, "fatal", color.Bit24(128, 0, 0, false).Sprintf},
		{logLevels.Error, "error", color.Bit24(255, 0, 0, false).Sprintf},
		{logLevels.Check, "check", color.Bit24(255, 255, 0, false).Sprintf},
		{logLevels.Warn, "warn ", color.Bit24(0, 255, 0, false).Sprintf},
		{logLevels.Info, "info ", color.Bit24(255, 255, 0, false).Sprintf},
		{logLevels.Debug, "debug", color.Bit24(0, 128, 255, false).Sprintf},
		{logLevels.Trace, "trace", color.Bit24(128, 0, 255, false).Sprintf},
	}
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
	LogChanDisabled = uberatomic.NewBool(true)
	LogChan         chan Entry
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

// AddLogChan adds a channel that log entries are sent to
func AddLogChan() (ch chan Entry) {
	LogChanDisabled.Store(false)
	if LogChan != nil {
		panic("warning warning")
	}
	// L.Writer.Write.Store( false
	LogChan = make(chan Entry)
	return LogChan
}

// GetLogPrinterSet returns a set of LevelPrinter with their subsystem preloaded
func GetLogPrinterSet(subsystem string) (Fatal, Error, Warn, Info, Debug, Trace LevelPrinter) {
	return _getOnePrinter(_Fatal, subsystem),
		_getOnePrinter(_Error, subsystem),
		_getOnePrinter(_Warn, subsystem),
		_getOnePrinter(_Info, subsystem),
		_getOnePrinter(_Debug, subsystem),
		_getOnePrinter(_Trace, subsystem)
}

func _getOnePrinter(level int32, subsystem string) LevelPrinter {
	return LevelPrinter{
		Ln:  _ln(level, subsystem),
		F:   _f(level, subsystem),
		S:   _s(level, subsystem),
		C:   _c(level, subsystem),
		Chk: _chk(level, subsystem),
	}
}

// SetLogLevel sets the log level via a string, which can be truncated down to
// one character, similar to nmcli's argument processor, as the first letter is
// unique. This could be used with a linter to make larger command sets.
func SetLogLevel(l string) {
	if l == "" {
		l = "info"
	}
	// fmt.Fprintln(os.Stderr, "setting log level", l)
	lvl := logLevels.Info
	for i := range LevelSpecs {
		if LevelSpecs[i].Name[:1] == l[:1] {
			lvl = LevelSpecs[i].ID
		}
	}
	currentLevel.Store(lvl)
}

// SetLogWriter atomically changes the log io.Writer interface
func SetLogWriter(wr io.Writer) {
	// w := unsafe.Pointer(writer)
	// c := unsafe.Pointer(wr)
	// atomic.SwapPointer(&w, c)
	writer = wr
}

func SetLogWriteToFile(path, appName string) (e error) {
	// copy existing log file to dated log file as we will truncate it per
	// session
	path = filepath.Join(path, "log"+appName)
	if _, e = os.Stat(path); e == nil {
		var b []byte
		b, e = ioutil.ReadFile(path)
		if e == nil {
			ioutil.WriteFile(path+fmt.Sprint(time.Now().Unix()), b, 0600)
		}
	}
	var fileWriter *os.File
	if fileWriter, e = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0600); e != nil {
		fmt.Fprintln(os.Stderr, "unable to write log to", path, "error:", e)
		return
	}
	mw := io.MultiWriter(os.Stderr, fileWriter)
	fileWriter.Write([]byte("logging to file '" + path + "'\n"))
	mw.Write([]byte("logging to file '" + path + "'\n"))
	SetLogWriter(mw)
	return
}

// SortSubsystemsList sorts the list of subsystems, to keep the data read-only,
// call this function right at the top of the main, which runs after
// declarations and main/init. Really this is just here to alert the reader.
func SortSubsystemsList() {
	sort.Strings(allSubsystems)
	// fmt.Fprintln(
	// 	os.Stderr,
	// 	spew.Sdump(allSubsystems),
	// 	spew.Sdump(highlighted),
	// 	spew.Sdump(logFilter),
	// )
}

// AddLoggerSubsystem adds a subsystem to the list of known subsystems and returns the
// string so it is nice and neat in the package logg.go file
func AddLoggerSubsystem(pathBase string) (subsystem string) {
	// var split []string
	var ok bool
	var file string
	_, file, _, ok = runtime.Caller(1)
	if ok {
		r := strings.Split(file, pathBase)
		// fmt.Fprintln(os.Stderr, version.PathBase, r)
		fromRoot := filepath.Base(file)
		if len(r) > 1 {
			fromRoot = r[1]
		}
		split := strings.Split(fromRoot, "/")
		// fmt.Fprintln(os.Stderr, version.PathBase, "file", file, r, fromRoot, split)
		subsystem = strings.Join(split[:len(split)-1], "/")
		// fmt.Fprintln(os.Stderr, "adding subsystem", subsystem)
		allSubsystems = append(allSubsystems, subsystem)
	}
	return
}

// StoreHighlightedSubsystems sets the list of subsystems to highlight
func StoreHighlightedSubsystems(highlights []string) (found bool) {
	highlightMx.Lock()
	highlighted = make(map[string]struct{}, len(highlights))
	for i := range highlights {
		highlighted[highlights[i]] = struct{}{}
	}
	highlightMx.Unlock()
	return
}

// LoadHighlightedSubsystems returns a copy of the map of highlighted subsystems
func LoadHighlightedSubsystems() (o []string) {
	highlightMx.Lock()
	o = make([]string, len(logFilter))
	var counter int
	for i := range logFilter {
		o[counter] = i
		counter++
	}
	highlightMx.Unlock()
	sort.Strings(o)
	return
}

// StoreSubsystemFilter sets the list of subsystems to filter
func StoreSubsystemFilter(filter []string) {
	_logFilterMx.Lock()
	logFilter = make(map[string]struct{}, len(filter))
	for i := range filter {
		logFilter[filter[i]] = struct{}{}
	}
	_logFilterMx.Unlock()
}

// LoadSubsystemFilter returns a copy of the map of filtered subsystems
func LoadSubsystemFilter() (o []string) {
	_logFilterMx.Lock()
	o = make([]string, len(logFilter))
	var counter int
	for i := range logFilter {
		o[counter] = i
		counter++
	}
	_logFilterMx.Unlock()
	sort.Strings(o)
	return
}

// _isHighlighted returns true if the subsystem is in the list to have attention
// getters added to them
func _isHighlighted(subsystem string) (found bool) {
	highlightMx.Lock()
	_, found = highlighted[subsystem]
	highlightMx.Unlock()
	return
}

// AddHighlightedSubsystem adds a new subsystem Name to the highlighted list
func AddHighlightedSubsystem(hl string) struct{} {
	highlightMx.Lock()
	highlighted[hl] = struct{}{}
	highlightMx.Unlock()
	return struct{}{}
}

// _isSubsystemFiltered returns true if the subsystem should not pr logs
func _isSubsystemFiltered(subsystem string) (found bool) {
	_logFilterMx.Lock()
	_, found = logFilter[subsystem]
	_logFilterMx.Unlock()
	return
}

// AddFilteredSubsystem adds a new subsystem Name to the highlighted list
func AddFilteredSubsystem(hl string) struct{} {
	_logFilterMx.Lock()
	logFilter[hl] = struct{}{}
	_logFilterMx.Unlock()
	return struct{}{}
}

func getTimeText(level int32) string {
	// since := time.Now().Sub(logger_started).Round(time.Millisecond).String()
	// diff := 12 - len(since)
	// if diff > 0 {
	// 	since = strings.Repeat(" ", diff) + since + " "
	// }
	return color.Bit24(99, 99, 99, false).Sprint(time.Now().
		Format(time.StampMilli))
}

func _ln(level int32, subsystem string) func(a ...interface{}) {
	return func(a ...interface{}) {
		if level <= currentLevel.Load() && !_isSubsystemFiltered(subsystem) {
			printer := fmt.Sprintf
			if _isHighlighted(subsystem) {
				printer = color.Bold.Sprintf
			}
			fmt.Fprintf(
				writer,
				printer(
					"%-58v%s%s%-6v %s\n",
					getLoc(2, level, subsystem),
					getTimeText(level),
					color.Bit24(20, 20, 20, true).
						Sprint(AppColorizer(" "+App)),
					LevelSpecs[level].Colorizer(
						color.Bit24(20, 20, 20, true).
							Sprint(" "+LevelSpecs[level].Name+" "),
					),
					AppColorizer(joinStrings(" ", a...)),
				),
			)
		}
	}
}

func _f(level int32, subsystem string) func(format string, a ...interface{}) {
	return func(format string, a ...interface{}) {
		if level <= currentLevel.Load() && !_isSubsystemFiltered(subsystem) {
			printer := fmt.Sprintf
			if _isHighlighted(subsystem) {
				printer = color.Bold.Sprintf
			}
			fmt.Fprintf(
				writer,
				printer(
					"%-58v%s%s%-6v %s\n",
					getLoc(2, level, subsystem),
					getTimeText(level),
					color.Bit24(20, 20, 20, true).
						Sprint(AppColorizer(" "+App)),
					LevelSpecs[level].Colorizer(
						color.Bit24(20, 20, 20, true).
							Sprint(" "+LevelSpecs[level].Name+" "),
					),
					AppColorizer(fmt.Sprintf(format, a...)),
				),
			)
		}
	}
}

func _s(level int32, subsystem string) func(a ...interface{}) {
	return func(a ...interface{}) {
		if level <= currentLevel.Load() && !_isSubsystemFiltered(subsystem) {
			printer := fmt.Sprintf
			if _isHighlighted(subsystem) {
				printer = color.Bold.Sprintf
			}
			fmt.Fprintf(
				writer,
				printer(
					"%-58v%s%s%s%s%s\n",
					getLoc(2, level, subsystem),
					getTimeText(level),
					color.Bit24(20, 20, 20, true).
						Sprint(AppColorizer(" "+App)),
					LevelSpecs[level].Colorizer(
						color.Bit24(20, 20, 20, true).
							Sprint(" "+LevelSpecs[level].Name+" "),
					),
					AppColorizer(
						" spew:",
					),
					fmt.Sprint(
						color.Bit24(20, 20, 20, true).Sprint("\n\n"+spew.Sdump(a)),
						"\n",
					),
				),
			)
		}
	}
}

func _c(level int32, subsystem string) func(closure func() string) {
	return func(closure func() string) {
		if level <= currentLevel.Load() && !_isSubsystemFiltered(subsystem) {
			printer := fmt.Sprintf
			if _isHighlighted(subsystem) {
				printer = color.Bold.Sprintf
			}
			fmt.Fprintf(
				writer,
				printer(
					"%-58v%s%s%-6v %s\n",
					getLoc(2, level, subsystem),
					getTimeText(level),
					color.Bit24(20, 20, 20, true).
						Sprint(AppColorizer(" "+App)),
					LevelSpecs[level].Colorizer(
						color.Bit24(20, 20, 20, true).
							Sprint(" "+LevelSpecs[level].Name+" "),
					),
					AppColorizer(closure()),
				),
			)
		}
	}
}

func _chk(level int32, subsystem string) func(e error) bool {
	return func(e error) bool {
		if level <= currentLevel.Load() && !_isSubsystemFiltered(subsystem) {
			if e != nil {
				printer := fmt.Sprintf
				if _isHighlighted(subsystem) {
					printer = color.Bold.Sprintf
				}
				fmt.Fprintf(
					writer,
					printer(
						"%-58v%s%s%-6v %s\n",
						getLoc(2, level, subsystem),
						getTimeText(level),
						color.Bit24(20, 20, 20, true).
							Sprint(AppColorizer(" "+App)),
						LevelSpecs[level].Colorizer(
							color.Bit24(20, 20, 20, true).
								Sprint(" "+LevelSpecs[level].Name+" "),
						),
						LevelSpecs[level].Colorizer(joinStrings(" ", e.Error())),
					),
				)
				return true
			}
		}
		return false
	}
}

// joinStrings constructs a string from an slice of interface same as Println but
// without the terminal newline
func joinStrings(sep string, a ...interface{}) (o string) {
	for i := range a {
		o += fmt.Sprint(a[i])
		if i < len(a)-1 {
			o += sep
		}
	}
	return
}

// getLoc calls runtime.Caller and formats as expected by source code editors
// for terminal hyperlinks
//
// Regular expressions and the substitution texts to make these clickable in
// Tilix and other RE hyperlink configurable terminal emulators:
//
// This matches the shortened paths generated in this command and printed at
// the very beginning of the line as this logger prints:
//
// ^((([\/a-zA-Z@0-9-_.]+/)+([a-zA-Z@0-9-_.]+)):([0-9]+))
//
// 		goland --line $5 $GOPATH/src/github.com/p9c/matrjoska/$2
//
// I have used a shell variable there but tilix doesn't expand them,
// so put your GOPATH in manually, and obviously change the repo subpath.
//
//
// Change the path to use with another repository's logging output (
// someone with more time on their hands could probably come up with
// something, but frankly the custom links feature of Tilix has the absolute
// worst UX I have encountered since the 90s...
// Maybe in the future this library will be expanded with a tool that more
// intelligently sets the path, ie from CWD or other cleverness.
//
// This matches full paths anywhere on the commandline delimited by spaces:
//
// ([/](([\/a-zA-Z@0-9-_.]+/)+([a-zA-Z@0-9-_.]+)):([0-9]+))
//
// 		goland --line $5 /$2
//
// Adapt the invocation to open your preferred editor if it has the capability,
// the above is for Jetbrains Goland
//
func getLoc(skip int, level int32, subsystem string) (output string) {
	_, file, line, _ := runtime.Caller(skip)
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "getloc panic on subsystem", subsystem, file)
		}
	}()
	split := strings.Split(file, subsystem)
	if len(split) < 2 {
		output = fmt.Sprint(
			color.White.Sprint(subsystem),
			color.Gray.Sprint(
				file, ":", line,
			),
		)
	} else {
		output = fmt.Sprint(
			color.White.Sprint(subsystem),
			color.Gray.Sprint(
				split[1], ":", line,
			),
		)
	}
	return
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

func FileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	return e == nil
}

func Caller(comment string, skip int) string {
	_, file, line, _ := runtime.Caller(skip + 1)
	o := fmt.Sprintf("%s: %s:%d", comment, file, line)
	// L.Debug(o)
	return o
}
