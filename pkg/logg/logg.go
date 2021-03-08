package logg

import (
	"fmt"
	"github.com/fatih/color"
	uberatomic "go.uber.org/atomic"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
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
	// LevelPrinter defines a set of terminal printing primitives that output with extra
	// data, time, log logLevelList, and code location
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
	levelSpec struct {
		id        int32
		name      string
		colorizer func(format string, a ...interface{}) string
	}
)

var (
	logger_started = time.Now()
	levelSpecs     = []levelSpec{
		{logLevels.Off, "off", nil},
		{logLevels.Fatal, "fatal", color.RedString},
		{logLevels.Error, "error", color.YellowString},
		{logLevels.Check, "check", color.GreenString},
		{logLevels.Warn, "warn", color.HiWhiteString},
		{logLevels.Info, "info", color.WhiteString},
		{logLevels.Debug, "debug", color.BlueString},
		{logLevels.Trace, "trace", color.MagentaString},
	}
	// repositoryPath is the part of the filesystem path that contains the application
	// repository.
	// todo: can be injected using stroy - and is kinda crap but fix it llater
	repositoryPath = "/home/loki/src/github.com/p9c/pod/"
	// sep is just a convenient shortcut for this very longwinded expression
	sep          = string(os.PathSeparator)
	currentLevel = uberatomic.NewInt32(logLevels.Info)
	// writer can be swapped out for any io.*writer* that you want to use instead of
	// stdout.
	writer = &os.Stderr
	// allSubsystems stores all of the package subsystem names found in the current
	// application
	allSubsystems []string
	// highlighted is a text that helps visually distinguish a log entry by category
	highlighted = make(map[string]struct{})
	// logFilter specifies a set of packages that will not pr logs
	logFilter = make(map[string]struct{})
	// mutexes to prevent concurrent map accesses
	highlightMx, _logFilterMx sync.Mutex
	// logLevels is a shorthand access that minimises possible name collisions in the
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
)

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
		Chk: _ok(level, subsystem),
	}
}

// SetLogLevel sets the log level via a string, which can be truncated down to
// one character, similar to nmcli's argument processor, as the first letter is
// unique. This could be used with a linter to make larger command sets.
func SetLogLevel(l string) {
	lvl := logLevels.Info
	for i := range levelSpecs {
		if levelSpecs[i].name[:len(l)] == l {
			lvl = levelSpecs[i].id
		}
	}
	currentLevel.Store(lvl)
}

// SetLogWriter atomically changes the log io.Writer interface
func SetLogWriter(wr *io.Writer) {
	w := unsafe.Pointer(writer)
	c := unsafe.Pointer(wr)
	atomic.SwapPointer(&w, c)
}

// SortSubsystemsList sorts the list of subsystems, to keep the data read-only,
// call this function right at the top of the main, which runs after
// declarations and main/init. Really this is just here to alert the reader.
func SortSubsystemsList() {
	sort.Strings(allSubsystems)
}

// AddLoggerSubsystem adds a subsystem to the list of known subsystems and returns the
// string so it is nice and neat in the package logg.go file
func AddLoggerSubsystem() (subsystem string) {
	var pkgPath string
	var split []string
	var ok bool
	_, pkgPath, _, ok = runtime.Caller(1)
	_ = ok
	fmt.Println(pkgPath, repositoryPath)
	fromRoot := strings.Split(pkgPath, repositoryPath)[1]
	split = strings.Split(fromRoot, sep)
	subsystem = strings.Join(split[:len(split)-1], "/")
	allSubsystems = append(allSubsystems, subsystem)
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

// _addHighlightedSubsystem adds a new subsystem name to the highlighted list
func _addHighlightedSubsystem(hl string) {
	highlightMx.Lock()
	highlighted[hl] = struct{}{}
	highlightMx.Unlock()
}

// _isSubsystemFiltered returns true if the subsystem should not pr logs
func _isSubsystemFiltered(subsystem string) (found bool) {
	_logFilterMx.Lock()
	_, found = logFilter[subsystem]
	_logFilterMx.Unlock()
	return
}

func _ln(level int32, subsystem string) func(a ...interface{}) {
	return func(a ...interface{}) {
		if level >= currentLevel.Load() || !_isSubsystemFiltered(subsystem) {
			if _isHighlighted(subsystem) {
				subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
			} else {
				subsystem = ""
			}
			// set the runtime.Caller depth to appropriate for the caller
			depth := 3
			if level == logLevels.Check {
				depth = 4
			}
			fmt.Fprintln(
				*writer,
				levelSpecs[level].colorizer(
					"%v %s%s [%s] %s",
					time.Now().Truncate(time.Millisecond).Sub(logger_started),
					levelSpecs[level].name,
					subsystem,
					color.WhiteString(joinStrings(" ", a)),
					getLoc(depth),
				),
			)
		}
	}
}

func _f(level int32, subsystem string) func(format string, a ...interface{}) {
	return func(format string, a ...interface{}) {
		if level >= currentLevel.Load() || !_isSubsystemFiltered(subsystem) {
			if _isHighlighted(subsystem) {
				subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
			}
			fmt.Fprintf(
				*writer, format,
				levelSpecs[level].colorizer(
					"%v %s%s [%s] %s",
					time.Now().Truncate(time.Millisecond).Sub(logger_started),
					levelSpecs[level].name,
					subsystem,
					color.WhiteString(joinStrings(" ", a)),
					getLoc(3),
				),
			)
		}
	}
}

func _s(level int32, subsystem string) func(a ...interface{}) {
	return func(a ...interface{}) {
		if level >= currentLevel.Load() || !_isSubsystemFiltered(subsystem) {
			if _isHighlighted(subsystem) {
				subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
			}
			fmt.Fprintln(
				*writer,
				levelSpecs[level].colorizer(
					"%v %s%s\n%s\n%s",
					time.Now().Truncate(time.Millisecond).Sub(logger_started),
					levelSpecs[level].name,
					subsystem,
					color.WhiteString(joinStrings("\n", a)),
					getLoc(3),
				),
			)
		}
	}
}

func _c(level int32, subsystem string) func(closure func() string) {
	return func(closure func() string) {
		if level >= currentLevel.Load() || !_isSubsystemFiltered(subsystem) {
			if _isHighlighted(subsystem) {
				subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
			} else {
				subsystem = ""
			}
			_ln(level, subsystem)(closure())
		}
	}
}

func _ok(level int32, subsystem string) func(e error) bool {
	return func(e error) bool {
		if level >= currentLevel.Load() || !_isSubsystemFiltered(subsystem) {
			if e != nil  {
				if _isHighlighted(subsystem) {
					subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
				} else {
					subsystem = ""
				}
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
func getLoc(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", file, line)
}
