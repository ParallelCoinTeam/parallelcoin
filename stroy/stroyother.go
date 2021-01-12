// +build !windows

package main

import (
	"os/exec"
)

func WindowsExec(split []string) (out *exec.Cmd) {
	panic("this function should not be called not on windows")
	return
}
