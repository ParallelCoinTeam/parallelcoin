package main

import (
	"github.com/p9c/pod/pkg/logg"
)

var subsystem = logg.AddLoggerSubsystem()
var F, E, W, I, D, T logg.LevelPrinter = logg.GetLogPrinterSet(subsystem)

func init() {
	// to filter out this package, uncomment the following
	// var _ = logg.AddFilteredSubsystem(subsystem)
	
	// to highlight this package, uncomment the following
	// var _ = logg.AddHighlightedSubsystem(subsystem)
	
	// these are here to test whether they are working
	// F.Ln("F.Ln")
	// E.Ln("E.Ln")
	// W.Ln("W.Ln")
	// I.Ln("I.Ln")
	// D.Ln("D.Ln")
	// F.Ln("T.Ln")
	// F.F("%s", "F.F")
	// E.F("%s", "E.F")
	// W.F("%s", "W.F")
	// I.F("%s", "I.F")
	// D.F("%s", "D.F")
	// T.F("%s", "T.F")
	// F.C(func() string { return "F.C" })
	// E.C(func() string { return "E.C" })
	// W.C(func() string { return "W.C" })
	// I.C(func() string { return "I.C" })
	// D.C(func() string { return "D.C" })
	// T.C(func() string { return "T.C" })
	// F.C(func() string { return "F.C" })
	// E.Chk(errors.New("E.Chk"))
	// W.Chk(errors.New("W.Chk"))
	// I.Chk(errors.New("I.Chk"))
	// D.Chk(errors.New("D.Chk"))
	// T.Chk(errors.New("T.Chk"))
}
