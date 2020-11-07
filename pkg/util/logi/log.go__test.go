package logi

import (
	"errors"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	L.SetLevel("trace", true, "olt")
	L.Trace("testing")
	L.Debug("testing")
	fmt.Println("'", L.Check("'", errors.New("this is a test")))
	L.Check("", nil)
	L.Info("testing")
	L.Warn("testing")
	L.Error("testing")
	L.Fatal("testing")

}
