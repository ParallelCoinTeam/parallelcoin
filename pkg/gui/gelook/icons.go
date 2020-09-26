// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"image"
	"image/color"

	"github.com/p9c/pod/pkg/gui/gelook/ico"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func NewDuoUIIcons() (i map[string]*DuoUIIcon) {
	i = make(map[string]*DuoUIIcon)
	i["Checked"] = mustIcon(NewDuoUIIcon(icons.ToggleCheckBox))
	i["Unchecked"] = mustIcon(NewDuoUIIcon(icons.ToggleCheckBoxOutlineBlank))
	i["RadioChecked"] = mustIcon(NewDuoUIIcon(icons.ToggleRadioButtonChecked))
	i["RadioUnchecked"] = mustIcon(NewDuoUIIcon(icons.ToggleRadioButtonUnchecked))
	i["iconCancel"] = mustIcon(NewDuoUIIcon(icons.NavigationCancel))
	i["iconOK"] = mustIcon(NewDuoUIIcon(icons.NavigationCheck))
	i["iconClose"] = mustIcon(NewDuoUIIcon(icons.NavigationClose))
	i["foldIn"] = mustIcon(NewDuoUIIcon(icons.ContentRemove))
	i["minimize"] = mustIcon(NewDuoUIIcon(icons.NavigationExpandMore))
	i["zoom"] = mustIcon(NewDuoUIIcon(icons.NavigationExpandLess))
	i["logo"] = mustIcon(NewDuoUIIcon(ico.ParallelCoin))
	i["overviewIcon"] = mustIcon(NewDuoUIIcon(icons.ActionHome))
	i["sendIcon"] = mustIcon(NewDuoUIIcon(icons.ActionStarRate))
	i["receiveIcon"] = mustIcon(NewDuoUIIcon(icons.NavigationArrowDropDown))
	i["addressBookIcon"] = mustIcon(NewDuoUIIcon(icons.ActionBook))
	i["historyIcon"] = mustIcon(NewDuoUIIcon(icons.ActionHistory))
	i["closeIcon"] = mustIcon(NewDuoUIIcon(icons.NavigationClose))
	i["settingsIcon"] = mustIcon(NewDuoUIIcon(icons.ActionSettings))
	i["blocksIcon"] = mustIcon(NewDuoUIIcon(icons.ActionExplore))
	i["networkIcon"] = mustIcon(NewDuoUIIcon(icons.ActionFingerprint))
	i["traceIcon"] = mustIcon(NewDuoUIIcon(icons.ActionTrackChanges))
	// i["consoleIcon"] = mustIcon(NewDuoUIIcon(icons.ActionInput))
	i["helpIcon"] = mustIcon(NewDuoUIIcon(icons.NavigationArrowDropDown))
	i["counterPlusIcon"] = mustIcon(NewDuoUIIcon(icons.ImageExposurePlus1))
	i["counterMinusIcon"] = mustIcon(NewDuoUIIcon(icons.ImageExposureNeg1))
	i["CommunicationImportExport"] = mustIcon(NewDuoUIIcon(icons.CommunicationImportExport))
	i["NotificationNetworkCheck"] = mustIcon(NewDuoUIIcon(icons.NotificationNetworkCheck))
	i["NotificationSync"] = mustIcon(NewDuoUIIcon(icons.NotificationSync))
	i["NotificationSyncDisabled"] = mustIcon(NewDuoUIIcon(icons.NotificationSyncDisabled))
	i["NotificationSyncProblem"] = mustIcon(NewDuoUIIcon(icons.NotificationSyncProblem))
	i["NotificationVPNLock"] = mustIcon(NewDuoUIIcon(icons.NotificationVPNLock))
	i["network"] = mustIcon(NewDuoUIIcon(icons.NotificationWiFi))
	i["MapsLayers"] = mustIcon(NewDuoUIIcon(icons.MapsLayers))
	i["MapsLayersClear"] = mustIcon(NewDuoUIIcon(icons.MapsLayersClear))
	i["ImageTimer"] = mustIcon(NewDuoUIIcon(icons.ImageTimer))
	i["ImageRemoveRedEye"] = mustIcon(NewDuoUIIcon(icons.ImageRemoveRedEye))
	i["DeviceSignalCellular0Bar"] = mustIcon(NewDuoUIIcon(icons.DeviceSignalCellular0Bar))
	i["DeviceWidgets"] = mustIcon(NewDuoUIIcon(icons.DeviceWidgets))
	i["ActionTimeline"] = mustIcon(NewDuoUIIcon(icons.ActionTimeline))
	i["HardwareWatch"] = mustIcon(NewDuoUIIcon(icons.HardwareWatch))
	i["consoleIcon"] = mustIcon(NewDuoUIIcon(icons.HardwareKeyboardHide))
	i["DeviceSignalCellular0Bar"] = mustIcon(NewDuoUIIcon(icons.DeviceSignalCellular0Bar))
	i["HardwareWatch"] = mustIcon(NewDuoUIIcon(icons.HardwareWatch))
	i["EditorMonetizationOn"] = mustIcon(NewDuoUIIcon(icons.EditorMonetizationOn))
	i["Run"] = mustIcon(NewDuoUIIcon(icons.AVPlayArrow))
	i["Stop"] = mustIcon(NewDuoUIIcon(icons.AVStop))
	i["Pause"] = mustIcon(NewDuoUIIcon(icons.AVPause))
	i["Kill"] = mustIcon(NewDuoUIIcon(icons.NavigationCancel))
	i["Restart"] = mustIcon(NewDuoUIIcon(icons.NavigationRefresh))
	i["Grab"] = mustIcon(NewDuoUIIcon(icons.NavigationMenu))
	i["Up"] = mustIcon(NewDuoUIIcon(icons.NavigationArrowDropUp))
	i["Down"] = mustIcon(NewDuoUIIcon(icons.NavigationArrowDropDown))
	i["iconGrab"] = mustIcon(NewDuoUIIcon(icons.NavigationMenu))
	i["iconUp"] = mustIcon(NewDuoUIIcon(icons.NavigationArrowDropUp))
	i["iconDown"] = mustIcon(NewDuoUIIcon(icons.NavigationArrowDropDown))
	i["Copy"] = mustIcon(NewDuoUIIcon(icons.ContentContentCopy))
	i["Paste"] = mustIcon(NewDuoUIIcon(icons.ContentContentPaste))
	i["Sidebar"] = mustIcon(NewDuoUIIcon(icons.ActionChromeReaderMode))
	i["Filter"] = mustIcon(NewDuoUIIcon(icons.ContentFilterList))
	i["FilterAll"] = mustIcon(NewDuoUIIcon(icons.ActionDoneAll))
	i["FilterNone"] = mustIcon(NewDuoUIIcon(icons.ContentBlock))
	i["Build"] = mustIcon(NewDuoUIIcon(icons.ActionBuild))
	i["Folded"] = mustIcon(NewDuoUIIcon(icons.NavigationChevronRight))
	i["Unfolded"] = mustIcon(NewDuoUIIcon(icons.NavigationExpandMore))
	i["HideAll"] = mustIcon(NewDuoUIIcon(icons.NavigationUnfoldLess))
	i["ShowAll"] = mustIcon(NewDuoUIIcon(icons.NavigationUnfoldMore))
	i["HideItem"] = mustIcon(NewDuoUIIcon(icons.ActionVisibilityOff))
	i["ShowItem"] = mustIcon(NewDuoUIIcon(icons.ActionVisibility))
	i["TRC"] = mustIcon(NewDuoUIIcon(icons.ActionSearch))
	i["DBG"] = mustIcon(NewDuoUIIcon(icons.ActionBugReport))
	i["INF"] = mustIcon(NewDuoUIIcon(icons.ActionInfo))
	i["WRN"] = mustIcon(NewDuoUIIcon(icons.ActionHelp))
	i["CHK"] = mustIcon(NewDuoUIIcon(icons.AlertWarning))
	i["ERR"] = mustIcon(NewDuoUIIcon(icons.AlertError))
	i["FTL"] = mustIcon(NewDuoUIIcon(icons.ImageFlashOn))
	i["Delete"] = mustIcon(NewDuoUIIcon(icons.ActionDelete))
	i["Send"] = mustIcon(NewDuoUIIcon(icons.ContentSend))
	return i
}

func mustIcon(ic *DuoUIIcon, err error) *DuoUIIcon {
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
