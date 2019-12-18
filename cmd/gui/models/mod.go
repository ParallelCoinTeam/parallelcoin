package models

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
)


type DuoUIcomponents struct {
	View               DuoUIcomponent
	Header             DuoUIcomponent
	Logo               DuoUIcomponent
	Body               DuoUIcomponent
	Sidebar            DuoUIcomponent
	Menu               DuoUIcomponent
	Content            DuoUIcomponent
	Overview           DuoUIcomponent
	OverviewTop        DuoUIcomponent
	SendReceive        DuoUIcomponent
	SendReceiveButtons DuoUIcomponent
	OverviewBottom     DuoUIcomponent
	Status             DuoUIcomponent
	History            DuoUIcomponent
	AddressBook        DuoUIcomponent
	Explorer           DuoUIcomponent
	Network            DuoUIcomponent
	Settings           DuoUIcomponent
}
type DuoUIcomponent struct {
	Layout layout.Flex
	Inset layout.Inset
}

type DuoUInav struct {
	Current       string
	IcoBackground color.RGBA
	IcoColor      color.RGBA
	IcoPadding    unit.Value
	IcoSize       unit.Value
	Overview      widget.Button
	History       widget.Button
	AddressBook   widget.Button
	Explorer      widget.Button
	Settings      widget.Button
}

type DuoUIicons struct {
	Logo        *material.Icon
	Overview    *material.Icon
	History     *material.Icon
	AddressBook *material.Icon
	Network     *material.Icon
	Explorer    *material.Icon
	Settings    *material.Icon
}

type Boot struct {
	IsBoot     bool `json:"boot"`
	IsFirstRun bool `json:"firstrun"`
	IsBootMenu bool `json:"menu"`
	IsBootLogo bool `json:"logo"`
	IsLoading  bool `json:"loading"`
}
