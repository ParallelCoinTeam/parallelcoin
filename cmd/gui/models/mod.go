package models

import (
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
	"github.com/p9c/pod/pkg/gio/widget/material"
	"image/color"
)

type DuoUIcomponents struct {
	View   DuoUIcomponent
	Header DuoUIcomponent
	//Intro              DuoUIcomponent
	Logo DuoUIcomponent
	//Log                DuoUIcomponent
	Body        DuoUIcomponent
	Sidebar     DuoUIcomponent
	Menu        DuoUIcomponent
	Content     DuoUIcomponent
	Overview    DuoUIcomponent
	OverviewTop DuoUIcomponent
	//SendReceive        DuoUIcomponent
	//SendReceiveButtons DuoUIcomponent
	OverviewBottom DuoUIcomponent
	Status         DuoUIcomponent
	StatusItem     DuoUIcomponent
	History        DuoUIcomponent
	AddressBook    DuoUIcomponent
	Explorer       DuoUIcomponent
	Network        DuoUIcomponent
	Console        DuoUIcomponent
	//ConsoleOutput      DuoUIcomponent
	//ConsoleInput       DuoUIcomponent
	Settings DuoUIcomponent
}

type DuoUIcomponent struct {
	Layout layout.Flex
	Inset  layout.Inset
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
	Console       widget.Button
	Settings      widget.Button
}

type DuoUIicons struct {
	Logo        *material.Icon
	Overview    *material.Icon
	History     *material.Icon
	AddressBook *material.Icon
	Network     *material.Icon
	Explorer    *material.Icon
	Console     *material.Icon
	Settings    *material.Icon
}

type Boot struct {
	IsBoot     bool `json:"boot"`
	IsFirstRun bool `json:"firstrun"`
	IsBootMenu bool `json:"menu"`
	IsBootLogo bool `json:"logo"`
	IsLoading  bool `json:"loading"`
}

type DuoUIconf struct {
	Abbrevation     string
	StatusTextColor color.RGBA
}
