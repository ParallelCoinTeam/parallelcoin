package gencerts

import (
	"github.com/p9c/pod/pkg/logg"
)

var subsystem = logg.AddLoggerSubsystem()
var ftl, err, wrn, inf, dbg, trc logg.LevelPrinter = logg.GetLogPrinterSet(subsystem)

func init() {
	// to filter out this package, uncomment the following
	// var _ = logg.AddFilteredSubsystem(subsystem)
	
	// to highlight this package, uncomment the following
	// var _ = logg.AddHighlightedSubsystem(subsystem)
	
	// these are here to test whether they are working
	// ftl.Ln("ftl.Ln")
	// err.Ln("err.Ln")
	// wrn.Ln("wrn.Ln")
	// inf.Ln("inf.Ln")
	// dbg.Ln("dbg.Ln")
	// trc.Ln("trc.Ln")
	// ftl.F("%s", "ftl.F")
	// err.F("%s", "err.F")
	// wrn.F("%s", "wrn.F")
	// inf.F("%s", "inf.F")
	// dbg.F("%s", "dbg.F")
	// trc.F("%s", "trc.F")
	// ftl.C(func() string { return "ftl.C" })
	// err.C(func() string { return "err.C" })
	// wrn.C(func() string { return "wrn.C" })
	// inf.C(func() string { return "inf.C" })
	// dbg.C(func() string { return "dbg.C" })
	// trc.C(func() string { return "trc.C" })
	// ftl.C(func() string { return "ftl.C" })
	// err.Chk(errors.New("err.Chk"))
	// wrn.Chk(errors.New("wrn.Chk"))
	// inf.Chk(errors.New("inf.Chk"))
	// dbg.Chk(errors.New("dbg.Chk"))
	// trc.Chk(errors.New("trc.Chk"))
}
