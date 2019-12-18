package duoui

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/cmd/gui/models"
)

type DuoUI struct {
	Boot *Boot
	rc   *RcVar
	cx   *conte.Xt
	ww   *app.Window
	gc   *layout.Context
	th   *material.Theme
	cs   *layout.Constraints
	ico  *models.DuoUIicons
	comp *models.DuoUIcomponents
	menu *models.DuoUInav
}
//
//type DuoUIcomponents struct {
//	view               DuoUIcomponent
//	header             DuoUIcomponent
//	logo               DuoUIcomponent
//	body               DuoUIcomponent
//	sidebar            DuoUIcomponent
//	menu               DuoUIcomponent
//	content            DuoUIcomponent
//	overview           DuoUIcomponent
//	overviewTop        DuoUIcomponent
//	sendReceive        DuoUIcomponent
//	sendReceiveButtons DuoUIcomponent
//	overviewBottom     DuoUIcomponent
//	status             DuoUIcomponent
//	history            DuoUIcomponent
//	addressbook        DuoUIcomponent
//	explorer           DuoUIcomponent
//	network            DuoUIcomponent
//	settings           DuoUIcomponent
//}
//type DuoUIcomponent struct {
//	l layout.Flex
//	i layout.Inset
//}
//
//type DuoUInav struct {
//	current       string
//	icoBackground color.RGBA
//	icoColor      color.RGBA
//	icoPadding    unit.Value
//	icoSize       unit.Value
//	overview      widget.Button
//	history       widget.Button
//	addressbook   widget.Button
//	explorer      widget.Button
//	settings      widget.Button
//}
//
//type DuoUIicons struct {
//	Logo        *material.Icon
//	Overview    *material.Icon
//	History     *material.Icon
//	AddressBook *material.Icon
//	Network     *material.Icon
//	Explorer    *material.Icon
//	Settings    *material.Icon
//}

type Boot struct {
	IsBoot     bool `json:"boot"`
	IsFirstRun bool `json:"firstrun"`
	IsBootMenu bool `json:"menu"`
	IsBootLogo bool `json:"logo"`
	IsLoading  bool `json:"loading"`
}
