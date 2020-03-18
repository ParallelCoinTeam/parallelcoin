package main

import (
	"time"

	"github.com/p9c/pod/pkg/logi"
)

func main() {
	logi.log.L.SetLevel("trace", true, "logi")
	for {
		logi.L.Trace("testing")
		// logi.L.Debug("testing")
		// fmt.Println("'", logi.L.Check(errors.New("this is a test")), "'")
		// logi.L.Check(nil)
		// logi.L.Info("testing")
		// logi.L.Warn("testing")
		// logi.L.Error("testing")
		// logi.L.Fatal("testing")
		time.Sleep(time.Second / 10)
	}

}
