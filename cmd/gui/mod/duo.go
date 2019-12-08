package mod

import (
	"gioui.org/app"
	"gioui.org/io/profile"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"net/http"
)

type DuoUI struct {
	Window    *app.Window
	Gtx       layout.Context
	Theme     material.Theme
	Layouts   DuoUIlayouts
	Ico       DuoUIicons
	Buttons   DuoUIbuttons
	usersList layout.List
	//users        []*user
	//userClicks   []gesture.Click
	//selectedUser *userPage
	//edit, edit2  *widget.Editor
	//fetchCommits func(u string)

	// Profiling.
	profiling   bool
	profile     profile.Event
	lastMallocs uint64

	Boot *Boot            `json:"boot"`
	Cf   *Configuration   `json:"cf"`
	Fs   *http.FileSystem `json:"fs"`
}

type Boot struct {
	IsBoot     bool `json:"boot"`
	IsFirstRun bool `json:"firstrun"`
	IsBootMenu bool `json:"menu"`
	IsBootLogo bool `json:"logo"`
	IsLoading  bool `json:"loading"`
}

type Configuration struct {
	Assets string `json:"assets"`
	Theme  bool   `json:"theme"`
	IsDev  bool   `json:"dev"`
}
