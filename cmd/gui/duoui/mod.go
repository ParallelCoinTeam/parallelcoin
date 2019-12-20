package duoui

import (
	"github.com/p9c/pod/pkg/gio/app"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/widget/material"
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
type Boot struct {
	IsBoot     bool `json:"boot"`
	IsFirstRun bool `json:"firstrun"`
	IsBootMenu bool `json:"menu"`
	IsBootLogo bool `json:"logo"`
	IsLoading  bool `json:"loading"`
}
