package main

import (
	"errors"
	"github.com/p9c/pod/pkg/logg"
	"os"
	"time"
	
	"github.com/p9c/pod/pkg/pipe/consume"
	"github.com/p9c/pod/pkg/util/qu"
)

func main() {
	// var e error
	logg.SetLogLevel("trace")
	// command := "pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := qu.T()
	// splitted := strings.Split(command, " ")
	splitted := os.Args[1:]
	w := consume.Log(quit, consume.SimpleLog(splitted[len(splitted)-1]), consume.FilterNone, splitted...)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> starting")
	consume.Start(w)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> started")
	time.Sleep(time.Second * 4)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopping")
	consume.Kill(w)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopped")
	// time.Sleep(time.Second * 5)
	// D.Ln(interrupt.GoroutineDump())
	// if e = w.Wait(); E.Chk(e) {
	// }
	// time.Sleep(time.Second * 3)
}

var subsystem = logg.AddLoggerSubsystem()
var ftl, err, wrn, inf, dbg, trc logg.LevelPrinter = logg.GetLogPrinterSet(subsystem)

func init() {
	// var _ = logg.AddFilteredSubsystem(subsystem)
	// var _ = logg.AddHighlightedSubsystem(subsystem)
	F.Ln("F.Ln")
	E.Ln("E.Ln")
	W.Ln("W.Ln")
	I.Ln("I.Ln")
	D.Ln("D.Ln")
	F.Ln("T.Ln")
	F.F("%s", "F.F")
	E.F("%s", "E.F")
	W.F("%s", "W.F")
	I.F("%s", "I.F")
	D.F("%s", "D.F")
	T.F("%s", "T.F")
	ftl.C(func() string { return "ftl.C" })
	err.C(func() string { return "err.C" })
	W.C(func() string { return "W.C" })
	I.C(func() string { return "inf.C" })
	D.C(func() string { return "D.C" })
	T.C(func() string { return "T.C" })
	ftl.C(func() string { return "ftl.C" })
	E.Chk(errors.New("E.Chk"))
	W.Chk(errors.New("W.Chk"))
	I.Chk(errors.New("inf.Chk"))
	D.Chk(errors.New("D.Chk"))
	T.Chk(errors.New("T.Chk"))
}
