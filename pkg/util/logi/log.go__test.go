package logi

import (
	"errors"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	L.SetLevel("trace", true, "olt")
	Trace("testing")
	Debug("testing")
	fmt.Println("'", Check(errors.New("this is a test")), "'")
	Check(nil)
	Info("testing")
	Warn("testing")
	Error("testing")
	Fatal("testing")

}
