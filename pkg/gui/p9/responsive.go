package p9

import (
	"fmt"

	"gioui.org/app"
	l "gioui.org/layout"
)

type WidgetMap map[string]l.Widget

var Resolutions = map[string]int{
	"nHd": 640,
	"qHD": 960,
	"HD":  1280,
	"FHD": 1920,
	"QHD": 2560,
	"UHD": 3840,
}

type Responsive struct {
	*Theme
	w       *app.Window
	widgets WidgetMap
}

func (th *Theme) Responsive(widgets WidgetMap) *Responsive {
	return &Responsive{
		Theme:   th,
		widgets: widgets,
	}
}

func (r *Responsive) Fn(gtx l.Context) l.Dimensions {
	width := gtx.Constraints.Max.X
	var out l.Widget
	switch {
	case width >= Resolutions["UHD"]:
		out = r.widgets["FHD"]
	case width >= Resolutions["QHD"]:
		out = r.widgets["QHD"]
	case width >= Resolutions["FHD"]:
		out = r.widgets["FHD"]
	case width >= Resolutions["HD"]:
		out = r.widgets["HD"]
	case width >= Resolutions["qHD"]:
		out = r.widgets["qHD"]
	case width >= Resolutions["nHD"]:
		out = r.widgets["nHD"]
	default: // this will never happen anyhow
		out = r.widgets["HD"]
	}
	return out(gtx)
}

func (r *Responsive) For(size string, widget l.Widget) *Responsive {
	var ok bool
	// stop programmer errors
	if _, ok = Resolutions[size]; !ok {
		panic(fmt.Sprintf("resolution does not exist: %s", size))
	}
	r.widgets[size] = widget
	return r
}
