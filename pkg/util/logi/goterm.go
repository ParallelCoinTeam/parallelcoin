//+build goterm

package logi

import "github.com/p9c/goterm"

func init() {
	TermWidth = goterm.Width
}
