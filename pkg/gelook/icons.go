// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"image"
	"image/color"

	"github.com/p9c/pod/pkg/gelook/ico"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func NewDuoUIicons() (i map[string]*DuoUIicon) {
	i = make(map[string]*DuoUIicon)
	i["Checked"] = mustIcon(NewDuoUIicon(icons.ToggleCheckBox))
	i["Unchecked"] = mustIcon(NewDuoUIicon(icons.ToggleCheckBoxOutlineBlank))
	i["RadioChecked"] = mustIcon(NewDuoUIicon(icons.ToggleRadioButtonChecked))
	i["RadioUnchecked"] = mustIcon(NewDuoUIicon(icons.ToggleRadioButtonUnchecked))
	i["iconCancel"] = mustIcon(NewDuoUIicon(icons.NavigationCancel))
	i["iconOK"] = mustIcon(NewDuoUIicon(icons.NavigationCheck))
	i["iconClose"] = mustIcon(NewDuoUIicon(icons.NavigationClose))
	i["foldIn"] = mustIcon(NewDuoUIicon(icons.ContentRemove))
	i["minimize"] = mustIcon(NewDuoUIicon(icons.NavigationExpandMore))
	i["zoom"] = mustIcon(NewDuoUIicon(icons.NavigationExpandLess))
	i["logo"] = mustIcon(NewDuoUIicon(ico.ParallelCoin))
	i["overviewIcon"] = mustIcon(NewDuoUIicon(icons.ActionHome))
	i["sendIcon"] = mustIcon(NewDuoUIicon(icons.ActionStarRate))
	i["receiveIcon"] = mustIcon(NewDuoUIicon(icons.NavigationArrowDropDown))
	i["addressBookIcon"] = mustIcon(NewDuoUIicon(icons.ActionBook))
	i["historyIcon"] = mustIcon(NewDuoUIicon(icons.ActionHistory))
	i["closeIcon"] = mustIcon(NewDuoUIicon(icons.NavigationClose))
	i["settingsIcon"] = mustIcon(NewDuoUIicon(icons.ActionSettings))
	i["blocksIcon"] = mustIcon(NewDuoUIicon(icons.ActionExplore))
	i["networkIcon"] = mustIcon(NewDuoUIicon(icons.ActionFingerprint))
	i["traceIcon"] = mustIcon(NewDuoUIicon(icons.ActionTrackChanges))
	// i["consoleIcon"] = mustIcon(NewDuoUIicon(icons.ActionInput))
	i["helpIcon"] = mustIcon(NewDuoUIicon(icons.NavigationArrowDropDown))
	i["counterPlusIcon"] = mustIcon(NewDuoUIicon(icons.ImageExposurePlus1))
	i["counterMinusIcon"] = mustIcon(NewDuoUIicon(icons.ImageExposureNeg1))
	i["CommunicationImportExport"] = mustIcon(NewDuoUIicon(icons.CommunicationImportExport))
	i["NotificationNetworkCheck"] = mustIcon(NewDuoUIicon(icons.NotificationNetworkCheck))
	i["NotificationSync"] = mustIcon(NewDuoUIicon(icons.NotificationSync))
	i["NotificationSyncDisabled"] = mustIcon(NewDuoUIicon(icons.NotificationSyncDisabled))
	i["NotificationSyncProblem"] = mustIcon(NewDuoUIicon(icons.NotificationSyncProblem))
	i["NotificationVPNLock"] = mustIcon(NewDuoUIicon(icons.NotificationVPNLock))
	i["network"] = mustIcon(NewDuoUIicon(icons.NotificationWiFi))
	i["MapsLayers"] = mustIcon(NewDuoUIicon(icons.MapsLayers))
	i["MapsLayersClear"] = mustIcon(NewDuoUIicon(icons.MapsLayersClear))
	i["ImageTimer"] = mustIcon(NewDuoUIicon(icons.ImageTimer))
	i["ImageRemoveRedEye"] = mustIcon(NewDuoUIicon(icons.ImageRemoveRedEye))
	i["DeviceSignalCellular0Bar"] = mustIcon(NewDuoUIicon(icons.DeviceSignalCellular0Bar))
	i["DeviceWidgets"] = mustIcon(NewDuoUIicon(icons.DeviceWidgets))
	i["ActionTimeline"] = mustIcon(NewDuoUIicon(icons.ActionTimeline))
	i["HardwareWatch"] = mustIcon(NewDuoUIicon(icons.HardwareWatch))
	i["consoleIcon"] = mustIcon(NewDuoUIicon(icons.HardwareKeyboardHide))
	i["DeviceSignalCellular0Bar"] = mustIcon(NewDuoUIicon(icons.DeviceSignalCellular0Bar))
	i["HardwareWatch"] = mustIcon(NewDuoUIicon(icons.HardwareWatch))
	i["EditorMonetizationOn"] = mustIcon(NewDuoUIicon(icons.EditorMonetizationOn))
	i["Run"] = mustIcon(NewDuoUIicon(icons.AVPlayArrow))
	i["Stop"] = mustIcon(NewDuoUIicon(icons.AVStop))
	i["Pause"] = mustIcon(NewDuoUIicon(icons.AVPause))
	i["Kill"] = mustIcon(NewDuoUIicon(icons.NavigationCancel))
	i["Restart"] = mustIcon(NewDuoUIicon(icons.NavigationRefresh))
	i["Grab"] = mustIcon(NewDuoUIicon(icons.NavigationMenu))
	i["Up"] = mustIcon(NewDuoUIicon(icons.NavigationArrowDropUp))
	i["Down"] = mustIcon(NewDuoUIicon(icons.NavigationArrowDropDown))
	i["iconGrab"] = mustIcon(NewDuoUIicon(icons.NavigationMenu))
	i["iconUp"] = mustIcon(NewDuoUIicon(icons.NavigationArrowDropUp))
	i["iconDown"] = mustIcon(NewDuoUIicon(icons.NavigationArrowDropDown))
	i["Copy"] = mustIcon(NewDuoUIicon(icons.ContentContentCopy))
	i["Paste"] = mustIcon(NewDuoUIicon(icons.ContentContentPaste))
	i["Sidebar"] = mustIcon(NewDuoUIicon(icons.ActionChromeReaderMode))
	i["Filter"] = mustIcon(NewDuoUIicon(icons.ContentFilterList))
	i["FilterAll"] = mustIcon(NewDuoUIicon(icons.ActionDoneAll))
	i["FilterNone"] = mustIcon(NewDuoUIicon(icons.ContentBlock))
	i["Build"] = mustIcon(NewDuoUIicon(icons.ActionBuild))
	i["Folded"] = mustIcon(NewDuoUIicon(icons.NavigationChevronRight))
	i["Unfolded"] = mustIcon(NewDuoUIicon(icons.NavigationExpandMore))
	i["HideAll"] = mustIcon(NewDuoUIicon(icons.NavigationUnfoldLess))
	i["ShowAll"] = mustIcon(NewDuoUIicon(icons.NavigationUnfoldMore))
	i["HideItem"] = mustIcon(NewDuoUIicon(icons.ActionVisibilityOff))
	i["ShowItem"] = mustIcon(NewDuoUIicon(icons.ActionVisibility))
	i["TRC"] = mustIcon(NewDuoUIicon(icons.ActionSearch))
	i["DBG"] = mustIcon(NewDuoUIicon(icons.ActionBugReport))
	i["INF"] = mustIcon(NewDuoUIicon(icons.ActionInfo))
	i["WRN"] = mustIcon(NewDuoUIicon(icons.ActionHelp))
	i["CHK"] = mustIcon(NewDuoUIicon(icons.AlertWarning))
	i["ERR"] = mustIcon(NewDuoUIicon(icons.AlertError))
	i["FTL"] = mustIcon(NewDuoUIicon(icons.ImageFlashOn))
	return i
}

func mustIcon(ic *DuoUIicon, err error) *DuoUIicon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
