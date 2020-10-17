// SPDX-License-Identifier: Unlicense OR MIT

package mod

import (
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/p9c/pod/pkg/gui/wallet/dap/res"
	"github.com/p9c/pod/pkg/gui/wallet/dap/win"
	"github.com/p9c/pod/pkg/gui/wallet/nav"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

//Duo App Platform
type Dap struct {
	Rc         *RcVar
	BeforeMain map[int]func()
	Main       func()
	AfterMain  map[int]func()
	//Ctx        context.Context
	//Tik        map[int]func()
	UI   *UserInterface
	S    *Settings
	Apps map[string]Sap
}

type Sap struct {
	Title string
	App   interface{}
}

type Settings struct {
	Dir  string
	File string
}

type UserInterface struct {
	Device string
	W      *win.Windows

	Theme *theme.Theme
	//Ekran   func(gtx C) D
	FontSize float32
	R        *res.Responsive
	//P        Pages
	N   *nav.Navigation
	F   func(gtx layout.Context) layout.Dimensions
	G   layout.Context
	Pop bool
	Ops op.Ops
}

//type Pages map[string]Page
