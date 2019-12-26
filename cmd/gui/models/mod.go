package models

import (
	"github.com/p9c/pod/pkg/gio/app"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
	"github.com/p9c/pod/pkg/gio/widget/material"
	"github.com/p9c/pod/pkg/pod"
	"image/color"
	"time"
)

type DuoUI struct {
	Boot    *Boot
	Ww      *app.Window
	Gc      *layout.Context
	Th      *material.Theme
	Cs      *layout.Constraints
	Ico     *DuoUIicons
	Comp    *DuoUIcomponents
	Menu    *DuoUInav
	Conf    *DuoUIconf
	Quit    chan struct{}
	Ready   chan struct{}
	IsReady bool
}

type DuoUIcomponents struct {
	View   DuoUIcomponent
	Header DuoUIcomponent
	// Intro              DuoUIcomponent
	Logo DuoUIcomponent
	// Log                DuoUIcomponent
	Body        DuoUIcomponent
	Sidebar     DuoUIcomponent
	Menu        DuoUIcomponent
	Content     DuoUIcomponent
	Overview    DuoUIcomponent
	OverviewTop DuoUIcomponent
	// SendReceive        DuoUIcomponent
	// SendReceiveButtons DuoUIcomponent
	OverviewBottom DuoUIcomponent
	Status         DuoUIcomponent
	StatusItem     DuoUIcomponent
	History        DuoUIcomponent
	AddressBook    DuoUIcomponent
	Explorer       DuoUIcomponent
	Network        DuoUIcomponent
	Console        DuoUIcomponent
	// ConsoleOutput      DuoUIcomponent
	// ConsoleInput       DuoUIcomponent
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
	IsBoot     bool   `json:"boot"`
	IsFirstRun bool   `json:"firstrun"`
	IsBootMenu bool   `json:"menu"`
	IsBootLogo bool   `json:"logo"`
	IsLoading  bool   `json:"loading"`
	IsScreen   string `json:"screen"`
}

type DuoUIconf struct {
	Abbrevation     string
	StatusTextColor color.RGBA
	Settings        DuoUIsettings
}

type DuoUIalert struct {
	Time      time.Time   `json:"time"`
	Title     string      `json:"title"`
	Message   interface{} `json:"message"`
	AlertType string      `json:"type"`
}

type DuoUIsettings struct {
	//db DuoUIdb
	//Display mod.DisplayConfig `json:"display"`
	Daemon DaemonConfig `json:"daemon"`
}

type DaemonConfig struct {
	Config *pod.Config `json:"config"`
	Schema pod.Schema  `json:"schema"`
}
