// SPDX-License-Identifier: Unlicense OR MIT

package view

import (
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/rcd"
)
type DuOS struct {
	Duo *model.DuoUI
	Rc *rcd.RcVar
	//Components map[string]*model.DuOScomponent
}

func DuOSboot()*DuOS{
	d := *new(DuOS)

	//duo := *new(model.DuoUI)
	//rc  := *new(rcd.RcVar)
	//
	//d.Duo = &duo
	//d.Rc  = &rc
	return &d
}
