package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/p9c/pod/pkg/util/logi"
)

func main() {
	logi.L.SetLevel("trace", true, "logi")
	for {
		trc.Ln("testing")
		logi.L.dbg.Ln("testing")
		fmt.Println("'", logi.L.Check("", errors.New("this is a test")))
		logi.L.Check("", nil)
		logi.L.inf.Ln("testing")
		logi.L.wrn.Ln("testing")
		logi.L.Error("testing")
		logi.L.ftl.Ln("testing")
		time.Sleep(time.Second / 10)
	}
}
