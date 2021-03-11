package main

import (
	"errors"
	"github.com/p9c/pod/pkg/logg"
	"os"
	"time"
	
	"github.com/p9c/pod/pkg/pipe/consume"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/qu"
)

func main() {
	// var e error
	logi.L.SetLevel("trace", false, "pod")
	// command := "pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := qu.T()
	// splitted := strings.Split(command, " ")
	splitted := os.Args[1:]
	w := consume.Log(quit, consume.SimpleLog(splitted[len(splitted)-1]), consume.FilterNone, splitted...)
	dbg.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> starting")
	consume.Start(w)
	dbg.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> started")
	time.Sleep(time.Second * 4)
	dbg.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopping")
	consume.Kill(w)
	dbg.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopped")
	// time.Sleep(time.Second * 5)
	// dbg.Ln(interrupt.GoroutineDump())
	// if e = w.Wait(); err.Chk(e) {
	// }
	// time.Sleep(time.Second * 3)
}

var subsystem = logg.AddLoggerSubsystem()
var ftl, err, wrn, inf, dbg, trc logg.LevelPrinter = logg.GetLogPrinterSet(subsystem)

func init() {
	// var _ = logg.AddFilteredSubsystem(subsystem)
	// var _ = logg.AddHighlightedSubsystem(subsystem)
	ftl.Ln("ftl.Ln")
	err.Ln("err.Ln")
	wrn.Ln("wrn.Ln")
	inf.Ln("inf.Ln")
	dbg.Ln("dbg.Ln")
	trc.Ln("trc.Ln")
	ftl.F("%s", "ftl.F")
	err.F("%s", "err.F")
	wrn.F("%s", "wrn.F")
	inf.F("%s", "inf.F")
	dbg.F("%s", "dbg.F")
	trc.F("%s", "trc.F")
	ftl.C(func() string { return "ftl.C" })
	err.C(func() string { return "err.C" })
	wrn.C(func() string { return "wrn.C" })
	inf.C(func() string { return "inf.C" })
	dbg.C(func() string { return "dbg.C" })
	trc.C(func() string { return "trc.C" })
	ftl.C(func() string { return "ftl.C" })
	err.Chk(errors.New("err.Chk"))
	wrn.Chk(errors.New("wrn.Chk"))
	inf.Chk(errors.New("inf.Chk"))
	dbg.Chk(errors.New("dbg.Chk"))
	trc.Chk(errors.New("trc.Chk"))
}
