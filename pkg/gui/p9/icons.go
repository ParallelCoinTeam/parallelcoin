// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"image/color"

	"github.com/p9c/pod/pkg/gui/ico"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

func NewIcons() (i map[string]*Icon) {
	i = make(map[string]*Icon)
	i["Checked"] = mustIcon(NewIcon(icons.ToggleCheckBox))
	i["Unchecked"] = mustIcon(NewIcon(icons.ToggleCheckBoxOutlineBlank))
	i["RadioChecked"] = mustIcon(NewIcon(icons.ToggleRadioButtonChecked))
	i["RadioUnchecked"] = mustIcon(NewIcon(icons.ToggleRadioButtonUnchecked))
	i["iconCancel"] = mustIcon(NewIcon(icons.NavigationCancel))
	i["iconOK"] = mustIcon(NewIcon(icons.NavigationCheck))
	i["iconClose"] = mustIcon(NewIcon(icons.NavigationClose))
	i["foldIn"] = mustIcon(NewIcon(icons.ContentRemove))
	i["minimize"] = mustIcon(NewIcon(icons.NavigationExpandMore))
	i["zoom"] = mustIcon(NewIcon(icons.NavigationExpandLess))
	i["logo"] = mustIcon(NewIcon(ico.ParallelCoin))
	i["overviewIcon"] = mustIcon(NewIcon(icons.ActionHome))
	i["sendIcon"] = mustIcon(NewIcon(icons.ActionStarRate))
	i["receiveIcon"] = mustIcon(NewIcon(icons.NavigationArrowDropDown))
	i["addressBookIcon"] = mustIcon(NewIcon(icons.ActionBook))
	i["historyIcon"] = mustIcon(NewIcon(icons.ActionHistory))
	i["closeIcon"] = mustIcon(NewIcon(icons.NavigationClose))
	i["settingsIcon"] = mustIcon(NewIcon(icons.ActionSettings))
	i["blocksIcon"] = mustIcon(NewIcon(icons.ActionExplore))
	i["networkIcon"] = mustIcon(NewIcon(icons.ActionFingerprint))
	i["traceIcon"] = mustIcon(NewIcon(icons.ActionTrackChanges))
	// i["consoleIcon"] = mustIcon(NewIcon(icons.ActionInput))
	i["helpIcon"] = mustIcon(NewIcon(icons.NavigationArrowDropDown))
	i["counterPlusIcon"] = mustIcon(NewIcon(icons.ImageExposurePlus1))
	i["counterMinusIcon"] = mustIcon(NewIcon(icons.ImageExposureNeg1))
	i["CommunicationImportExport"] = mustIcon(NewIcon(icons.CommunicationImportExport))
	i["NotificationNetworkCheck"] = mustIcon(NewIcon(icons.NotificationNetworkCheck))
	i["NotificationSync"] = mustIcon(NewIcon(icons.NotificationSync))
	i["NotificationSyncDisabled"] = mustIcon(NewIcon(icons.NotificationSyncDisabled))
	i["NotificationSyncProblem"] = mustIcon(NewIcon(icons.NotificationSyncProblem))
	i["NotificationVPNLock"] = mustIcon(NewIcon(icons.NotificationVPNLock))
	i["network"] = mustIcon(NewIcon(icons.NotificationWiFi))
	i["MapsLayers"] = mustIcon(NewIcon(icons.MapsLayers))
	i["MapsLayersClear"] = mustIcon(NewIcon(icons.MapsLayersClear))
	i["ImageTimer"] = mustIcon(NewIcon(icons.ImageTimer))
	i["ImageRemoveRedEye"] = mustIcon(NewIcon(icons.ImageRemoveRedEye))
	i["DeviceSignalCellular0Bar"] = mustIcon(NewIcon(icons.DeviceSignalCellular0Bar))
	i["DeviceWidgets"] = mustIcon(NewIcon(icons.DeviceWidgets))
	i["ActionTimeline"] = mustIcon(NewIcon(icons.ActionTimeline))
	i["HardwareWatch"] = mustIcon(NewIcon(icons.HardwareWatch))
	i["consoleIcon"] = mustIcon(NewIcon(icons.HardwareKeyboardHide))
	i["DeviceSignalCellular0Bar"] = mustIcon(NewIcon(icons.DeviceSignalCellular0Bar))
	i["HardwareWatch"] = mustIcon(NewIcon(icons.HardwareWatch))
	i["EditorMonetizationOn"] = mustIcon(NewIcon(icons.EditorMonetizationOn))
	i["Run"] = mustIcon(NewIcon(icons.AVPlayArrow))
	i["Stop"] = mustIcon(NewIcon(icons.AVStop))
	i["Pause"] = mustIcon(NewIcon(icons.AVPause))
	i["Kill"] = mustIcon(NewIcon(icons.NavigationCancel))
	i["Restart"] = mustIcon(NewIcon(icons.NavigationRefresh))
	i["Grab"] = mustIcon(NewIcon(icons.NavigationMenu))
	i["Up"] = mustIcon(NewIcon(icons.NavigationArrowDropUp))
	i["Down"] = mustIcon(NewIcon(icons.NavigationArrowDropDown))
	i["iconGrab"] = mustIcon(NewIcon(icons.NavigationMenu))
	i["iconUp"] = mustIcon(NewIcon(icons.NavigationArrowDropUp))
	i["iconDown"] = mustIcon(NewIcon(icons.NavigationArrowDropDown))
	i["Copy"] = mustIcon(NewIcon(icons.ContentContentCopy))
	i["Paste"] = mustIcon(NewIcon(icons.ContentContentPaste))
	i["Sidebar"] = mustIcon(NewIcon(icons.ActionChromeReaderMode))
	i["Filter"] = mustIcon(NewIcon(icons.ContentFilterList))
	i["FilterAll"] = mustIcon(NewIcon(icons.ActionDoneAll))
	i["FilterNone"] = mustIcon(NewIcon(icons.ContentBlock))
	i["Build"] = mustIcon(NewIcon(icons.ActionBuild))
	i["Folded"] = mustIcon(NewIcon(icons.NavigationChevronRight))
	i["Unfolded"] = mustIcon(NewIcon(icons.NavigationExpandMore))
	i["HideAll"] = mustIcon(NewIcon(icons.NavigationUnfoldLess))
	i["ShowAll"] = mustIcon(NewIcon(icons.NavigationUnfoldMore))
	i["HideItem"] = mustIcon(NewIcon(icons.ActionVisibilityOff))
	i["ShowItem"] = mustIcon(NewIcon(icons.ActionVisibility))
	i["TRC"] = mustIcon(NewIcon(icons.ActionSearch))
	i["DBG"] = mustIcon(NewIcon(icons.ActionBugReport))
	i["INF"] = mustIcon(NewIcon(icons.ActionInfo))
	i["WRN"] = mustIcon(NewIcon(icons.ActionHelp))
	i["CHK"] = mustIcon(NewIcon(icons.AlertWarning))
	i["ERR"] = mustIcon(NewIcon(icons.AlertError))
	i["FTL"] = mustIcon(NewIcon(icons.ImageFlashOn))
	i["Delete"] = mustIcon(NewIcon(icons.ActionDelete))
	i["Send"] = mustIcon(NewIcon(icons.ContentSend))
	return i
}

func mustIcon(ic *Icon, err error) *Icon {
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
//
// func fill(gtx *layout.Context, col color.RGBA) {
// 	cs := gtx.Constraints
// 	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
// 	dr := f32.Rectangle{
// 		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
// 	}
// 	paint.ColorOp{Color: col}.Add(gtx.Ops)
// 	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
// 	gtx.Dimensions = layout.Dimensions{Size: d}
// }
